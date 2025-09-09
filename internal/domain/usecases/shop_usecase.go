package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/dataSources/rabbitmq"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	dto "nusantara_service/internal/dto/request"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	"nusantara_service/internal/workers/payload"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ShopService struct {
	shopRepo    repositories.ShopRepository
	productRepo repositories.ProductRepository
	rdb         *redis.Client
	cld         cloudinary.CloudinaryService
}

func NewShopUsecase(repo repositories.ShopRepository, redis *redis.Client, cld *cloudinary.CloudinaryService, prodRepo repositories.ProductRepository) services.ShopService {
	return &ShopService{
		shopRepo:    repo,
		rdb:         redis,
		cld:         *cld,
		productRepo: prodRepo,
	}
}

func fileToBytes(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return io.ReadAll(f)
}

// CreateShop implements services.ShopService.
func (s *ShopService) CreateShop(ctx context.Context, superAdminId uuid.UUID, req dto.CreateShopRequest) error {
	user, err := s.shopRepo.FindUserIDSuperAdmin(ctx, superAdminId.String())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if user == nil {
		return response.NewCustomError(response.ErrForbidden, "user not permission", 403)
	}

	// Validasi
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return response.NewCustomError(response.ErrBadRequest, "name is required", 400)
	}
	if strings.TrimSpace(req.Description) == "" {
		return response.NewCustomError(response.ErrBadRequest, "description is required", 400)
	}
	if strings.TrimSpace(req.FullAddress) == "" {
		return response.NewCustomError(response.ErrBadRequest, "full_address is required", 400)
	}
	if req.Lat == nil {
		return response.NewCustomError(response.ErrBadRequest, "lat is required", 400)
	}
	if req.Lng == nil {
		return response.NewCustomError(response.ErrBadRequest, "lang is required", 400)
	}
	if req.Status == nil {
		return response.NewCustomError(response.ErrBadRequest, "status is required", 400)
	}
	if req.CashierIDs == nil {
		return response.NewCustomError(response.ErrBadRequest, "cashier is required", 400)
	}
	if req.Products == nil {
		return response.NewCustomError(response.ErrBadRequest, "product is required", 400)
	}
	if req.Cover == nil {
		return response.NewCustomError(response.ErrBadRequest, "cover is required", 400)
	}
	if req.Gallery == nil {
		return response.NewCustomError(response.ErrBadRequest, "gallery is required", 400)
	}

	// Create entity
	shop := &entities.ShopEntity{
		ID:          uuid.New(),
		Name:        name,
		Description: req.Description,
		FullAddress: req.FullAddress,
		Lat:         *req.Lat,
		Lng:         *req.Lng,
		Status:      *req.Status,
		CreatedBy:   uuid.MustParse(user.ID),
	}
	if err := s.shopRepo.Create(ctx, shop); err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to create shop", 500)
	}

	// Upload cover & gallery pakai RabbitMQ
	if req.Cover != nil {
		if bin, err := fileToBytes(req.Cover); err == nil && len(bin) > 0 {
			task := payload.ImageUploadPayload{
				ShopID:    shop.ID,
				Type:      "cover",
				FileBytes: bin,
				Folder:    "nusantara_service/shops",
				Filename:  fmt.Sprintf("shop_%s_cover", shop.ID.String()),
			}
			_ = rabbitmq.PublishToQueue("", rabbitmq.SendImageShopQueueName, task)
		}
	}
	for i, g := range req.Gallery {
		if g == nil {
			continue
		}
		if bin, err := fileToBytes(g); err == nil && len(bin) > 0 {
			task := payload.ImageUploadPayload{
				ShopID:    shop.ID,
				Type:      "gallery",
				FileBytes: bin,
				Folder:    "nusantara_service/shops",
				Filename:  fmt.Sprintf("shops_%s_gallery_%d", shop.ID.String(), i+1),
			}
			_ = rabbitmq.PublishToQueue("", rabbitmq.SendImageShopQueueName, task)
		}
	}

	// Assign Cashiers
	if len(req.CashierIDs) > 0 {
		for _, cashierId := range req.CashierIDs {
			if _, err := s.shopRepo.FindCashier(ctx, cashierId); err != nil {
				return response.NewCustomError(response.ErrBadRequest, fmt.Sprintf("invalid cashier ID: %s", cashierId), 400)
			}
		}
		if err := s.shopRepo.AssignCashiers(ctx, shop.ID, req.CashierIDs); err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to assign cashier", 500)
		}
	}

	// Assign Products
	var shopProducts []entities.ShopProductEntity
	for _, p := range req.Products {
		master, err := s.productRepo.GetByID(ctx, p.ProductID)
		if err != nil {
			return response.NewCustomError(response.ErrNotFound, "product not found", 404)
		}

		price := master.Price
		if p.Price != nil {
			price = int(*p.Price)
		}
		status := master.Status
		if p.Status != nil {
			status = *p.Status
		}

		shopProducts = append(shopProducts, entities.ShopProductEntity{
			ID:        uuid.New(),
			ShopID:    shop.ID,
			ProductID: p.ProductID,
			Price:     float64(price),
			Status:    status,
			Stock:     p.Stock,
		})
	}

	if len(shopProducts) > 0 {
		if err := s.shopRepo.AssignProducts(ctx, shop.ID, shopProducts); err != nil {
			return err
		}
	}

	s.invalidateShopCache(ctx)
	return nil
}

// GetAllShop implements services.ShopService.
func (s *ShopService) GetAllShop(ctx context.Context, page int, limit int, search string) ([]*entities.ShopEntity, int, error) {
	cacheKey := fmt.Sprintf("shops:search:%s:page:%d:limit:%d", search, page, limit)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*entities.ShopEntity `json:"data"`
			Total int                    `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit

	items, total, err := s.shopRepo.FindAll(ctx, offset, limit, search)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, response.NewCustomError(response.ErrNotFound, "Shop Not found", 400)
		}
		return nil, 0, err
	}

	buf, _ := json.Marshal(map[string]any{
		"data":  items,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return items, total, nil
}

// GetByIdShop implements services.ShopService.
func (s *ShopService) GetByIdShop(ctx context.Context, id uuid.UUID) (*entities.ShopEntity, error) {
	cacheKey := fmt.Sprintf("shop:%s", id)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var shop entities.ShopEntity
		if json.Unmarshal([]byte(cached), &shop) == nil {
			return &shop, nil
		}
	}

	shop, err := s.shopRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "shop not found", 404)
		}
		return nil, err
	}

	buf, _ := json.Marshal(shop)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return shop, nil
}

// UpdateShop implements services.ShopService.
func (s *ShopService) UpdateShop(ctx context.Context, superAdminId uuid.UUID, id uuid.UUID, req dto.UpdateShopRequest) error {
	shop, err := s.shopRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "shop not found", 404)
		}
		return response.NewCustomError(response.ErrInternal, "failed to find shop", 500)
	}

	// Update basic fields
	if req.Name != "" {
		shop.Name = req.Name
	}
	if req.Description != "" {
		shop.Description = req.Description
	}
	if req.FullAddress != "" {
		shop.FullAddress = req.FullAddress
	}
	if req.Lat != 0 {
		shop.Lat = req.Lat
	}
	if req.Lng != 0 {
		shop.Lng = req.Lng
	}

	// Cover update
	if req.NewCover != nil {
		if shop.Cover != "" {
			publicID := utils.ExtractPublicIDFromCloudinaryURL(shop.Cover)
			if publicID != "" {
				_ = s.cld.DestroyImage(ctx, publicID)
			}
		}
		if bin, err := fileToBytes(req.NewCover); err == nil && len(bin) > 0 {
			task := payload.ImageUploadPayload{
				ShopID:    shop.ID,
				Type:      "cover",
				FileBytes: bin,
				Folder:    "nusantara_service/shops",
				Filename:  fmt.Sprintf("shop_%s_cover", shop.ID.String()),
			}
			_ = rabbitmq.PublishToQueue("", rabbitmq.SendImageShopQueueName, task)
		}
	}

	// Gallery update
	if len(req.NewGallery) > 0 {
		if req.ReplaceGallery {
			_ = s.shopRepo.DeleteGalleryByShopID(ctx, shop.ID)
			for _, gi := range shop.ShopImages {
				if gi.Image.ImagePath != "" {
					publicID := utils.ExtractPublicIDFromCloudinaryURL(gi.Image.ImagePath)
					if publicID != "" {
						_ = s.cld.DestroyImage(ctx, publicID)
					}
				}
			}
		}
		for i, g := range req.NewGallery {
			if bin, err := fileToBytes(g); err == nil && len(bin) > 0 {
				task := payload.ImageUploadPayload{
					ShopID:    shop.ID,
					Type:      "gallery",
					FileBytes: bin,
					Folder:    "nusantara_service/shops",
					Filename:  fmt.Sprintf("shops_%s_gallery_%d", shop.ID.String(), i+1),
				}
				_ = rabbitmq.PublishToQueue("", rabbitmq.SendImageShopQueueName, task)
			}
		}
	}

	// Replace Cashiers
	if len(req.CashierIDs) > 0 {
		if err := s.shopRepo.ReplaceCashiers(ctx, shop.ID, req.CashierIDs); err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to update cashiers", 500)
		}
	}

	// Replace Products
	if len(req.Products) > 0 {
		_ = s.shopRepo.DeleteShopProductsByShopID(ctx, shop.ID)

		var shopProducts []entities.ShopProductEntity
		for _, p := range req.Products {
			master, err := s.productRepo.GetByID(ctx, p.ProductID)
			if err != nil {
				return response.NewCustomError(response.ErrNotFound, "product not found", 404)
			}

			price := master.Price
			if p.Price != nil {
				price = int(*p.Price)
			}
			status := master.Status
			if p.Status != nil {
				status = *p.Status
			}

			shopProducts = append(shopProducts, entities.ShopProductEntity{
				ID:        uuid.New(),
				ShopID:    shop.ID,
				ProductID: p.ProductID,
				Price:     float64(price),
				Status:    status,
				Stock:     p.Stock,
			})
		}
		if len(shopProducts) > 0 {
			if err := s.shopRepo.ReplaceProducts(ctx, shop.ID, shopProducts); err != nil {
				return err
			}
		}
	}

	if err := s.shopRepo.Update(ctx, id, shop); err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to update shop", 500)
	}

	s.invalidateShopCache(ctx)
	return nil
}

// Delete implements services.ShopService.
func (s *ShopService) Delete(ctx context.Context, id uuid.UUID) error {
	shop, err := s.shopRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "shop not found", 404)
		}
		return response.NewCustomError(response.ErrInternal, "failed to find shop", 500)
	}

	if len(shop.ShopImages) > 0 {
		_ = s.shopRepo.DeleteGalleryByShopID(ctx, shop.ID)
	}

	if shop.Cover != "" {
		publicID := utils.ExtractPublicIDFromCloudinaryURL(shop.Cover)
		if publicID != "" {
			if err := s.cld.DestroyImage(ctx, publicID); err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to delete image", 500)
			}
		}
	}

	err = s.shopRepo.Delete(ctx, shop.ID)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to delete shop", 500)
	}

	s.invalidateShopCache(ctx)
	return nil
}

// UpdateStatusShop implements services.ShopService.
func (s *ShopService) UpdateStatusShop(ctx context.Context, id uuid.UUID, req dto.UpdateStatusShopRequest) error {
	if req.Status < 0 || req.Status > 1 {
		return response.NewCustomError(response.ErrBadRequest, "status must be 0 or 1", 400)
	}

	if err := s.shopRepo.UpdateStatus(ctx, id, req.Status); err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to update status", 500)
	}

	s.invalidateShopCache(ctx)
	return nil
}

func (s *ShopService) invalidateShopCache(ctx context.Context) {
	iter := s.rdb.Scan(ctx, 0, "shops:*", 0).Iterator()
	for iter.Next(ctx) {
		s.rdb.Del(ctx, iter.Val())
	}

	iterID := s.rdb.Scan(ctx, 0, "shop:*", 0).Iterator()
	for iterID.Next(ctx) {
		s.rdb.Del(ctx, iterID.Val())
	}
}
