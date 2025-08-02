package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TypeProductService struct {
	repo       repositories.TypeProductRepository
	rdb        *redis.Client
	cloudinary cloudinary.CloudinaryService
}

func NewTypeProductUsecase(repo repositories.TypeProductRepository, rdb *redis.Client, clodinary *cloudinary.CloudinaryService) services.TypeProductService {
	return &TypeProductService{repo: repo, rdb: rdb, cloudinary: *clodinary}
}

// CreateTypeProduct implements services.TypeProductService.
func (t *TypeProductService) CreateTypeProduct(ctx context.Context, userId string, req dto.CreateTypeProductRequest, image *multipart.FileHeader) (*entities.TypeProductEntity, error) {
	existing, err := t.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, response.NewCustomError(response.ErrNotFound, err.Error(), 404)
	}

	name := strings.TrimSpace(req.Name)

	if existing != nil {
		return nil, response.NewCustomError(response.ErrExists, "name alredy exists", 409)
	}

	if name == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "name is required", 400)
	}

	if image == nil {
		return nil, response.NewCustomError(response.ErrBadRequest, "image is required", 400)
	}

	if req.Status < 0 || req.Status > 1 {
		return nil, response.NewCustomError(response.ErrBadRequest, "Status must be 0 & 1", 400)
	}

	user, err := t.repo.FindByUserIDSuperAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	folder := "nusantara_service/type_product"
	filename := fmt.Sprintf("type_product_%s", req.Name)
	imageUrl, err := t.cloudinary.UploadImage(ctx, image, folder, filename)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to upload image", 500)
	}

	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, response.NewCustomError(response.ErrBadRequest, "invalid user ID", 400)
	}

	newTypeProduct := &entities.TypeProductEntity{
		Name:   name,
		Status: req.Status,
		Image:  imageUrl,
		UserID: userUUID,
		User:   *user,
	}

	_, err = t.repo.Create(ctx, newTypeProduct)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to save type product", 500)
	}

	t.InvalidateTypeProductCache(ctx)
	return newTypeProduct, nil
}

// GetAllTypeProduct implements services.TypeProductService.
func (t *TypeProductService) GetAllTypeProduct(ctx context.Context, page int, limit int, search string) ([]*entities.TypeProductEntity, int, error) {
	cacheKey := fmt.Sprintf("type_products:search:%s:page:%d:limit:%d", search, page, limit)
	cached, err := configs.GetRedis(ctx, cacheKey)
	if err == nil {
		var result struct {
			Data  []*entities.TypeProductEntity `json:"data"`
			Total int                           `json:"total"`
		}

		_ = json.Unmarshal([]byte(cached), &result)
		return result.Data, result.Total, nil
	}

	offset := (page - 1) * limit
	// rediskey := fmt.Sprintf("type_products:page=%d&limit=%d", page, limit)

	// cached, err := configs.GetRedis(ctx, rediskey)
	// if err == nil && cached != "" {
	// 	var cachedData struct {
	// 		TypeProducts []*entities.TypeProductEntity `json:"type_products"`
	// 		Total        int                           `json:"total"`
	// 	}
	// 	if err := json.Unmarshal([]byte(cached), &cachedData); err != nil {
	// 		return cachedData.TypeProducts, cachedData.Total, nil
	// 	}
	// }

	typeProducts, err := t.repo.FindAllWithSearch(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrNotFound, "Type Product Not Found", 404)
	}

	total, err := t.repo.CountAllWithSearch(ctx, search)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrBadRequest, "failed to count type product", 400)
	}
	dataToCache, _ := json.Marshal(map[string]interface{}{
		"data":  typeProducts,
		"total": total,
	})

	_ = configs.SetRedis(ctx, cacheKey, dataToCache, time.Minute*10)
	return typeProducts, total, nil
}

// GetByIdTypeProduct implements services.TypeProductService.
func (t *TypeProductService) GetByIdTypeProduct(ctx context.Context, id uuid.UUID) (*entities.TypeProductEntity, error) {
	redisKey := fmt.Sprintf("type_product:%s", id.String())

	cached, err := configs.GetRedis(ctx, redisKey)
	if err == nil && cached != "" {
		var typeProduct entities.TypeProductEntity
		if err := json.Unmarshal([]byte(cached), &typeProduct); err != nil {
			return &typeProduct, nil
		}
	}

	typeProduct, err := t.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "Not Found Type Product", 404)
		}
		return nil, err
	}

	dataToCache, _ := json.Marshal(typeProduct)
	_ = configs.SetRedis(ctx, redisKey, dataToCache, time.Minute*30)

	return typeProduct, nil
}

// UpdateTypeProduct implements services.TypeProductService.
func (t *TypeProductService) UpdateTypeProduct(ctx context.Context, userId string, id uuid.UUID, req dto.UpdateTypeProductRequest, image *multipart.FileHeader) (*entities.TypeProductEntity, error) {
	typeProduct, err := t.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "Not Found Type Product", 404)
		}
		return nil, err
	}

	name := strings.TrimSpace(req.Name)

	if name != "" {
		existing, _ := t.repo.FindByName(ctx, req.Name)
		if existing != nil && existing.ID != id {
			return nil, response.NewCustomError(response.ErrExists, "Type Product already exists", 404)
		}
		typeProduct.Name = name
	}

	if image != nil {
		if typeProduct.Image != "" {
			publicID := utils.ExtractPublicIDFromCloudinaryURL(typeProduct.Image)
			if publicID != "" {
				if err := t.cloudinary.DestroyImage(ctx, publicID); err != nil {
					return nil, response.NewCustomError(response.ErrInternal, "failed to delete image", 500)
				}
			}
		}
		folder := "nusantara_service/type_product"
		filename := fmt.Sprintf("type_product_%s", id.String())
		imageUrl, err := t.cloudinary.UploadImage(ctx, image, folder, filename)
		if err != nil {
			return nil, response.NewCustomError(response.ErrInternal, "failed to upload image", 500)
		}

		typeProduct.Image = imageUrl

	}

	user, err := t.repo.FindByUserIDSuperAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, response.NewCustomError(response.ErrBadRequest, "invalid user ID", 400)
	}

	typeProduct.UserID = userUUID

	updatedTypeProduct, err := t.repo.Update(ctx, id, typeProduct)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "Failed to update type product", 500)
	}

	t.InvalidateTypeProductCache(ctx)

	return updatedTypeProduct, nil

}

// DeleteTypeProduct implements services.TypeProductService.
func (t *TypeProductService) DeleteTypeProduct(ctx context.Context, id uuid.UUID) error {
	typeProduct, err := t.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "Not Found Type Product", 404)
		}
		return err
	}

	if typeProduct.Image != "" {
		publicID := utils.ExtractPublicIDFromCloudinaryURL(typeProduct.Image)
		if publicID != "" {
			if err := t.cloudinary.DestroyImage(ctx, publicID); err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to delete image", 500)
			}
		}
	}

	t.InvalidateTypeProductCache(ctx)

	return t.repo.Delete(ctx, id)
}

// UpdateStatusTypeProduct implements services.TypeProductService.
func (t *TypeProductService) UpdateStatusTypeProduct(ctx context.Context, id uuid.UUID, req dto.UpdateStatusTypeProductRequest) error {
	typeProduct, err := t.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "Not Found Type Product", 404)
		}
		return err
	}
	if req.Status < 0 || req.Status > 1 {
		return response.NewCustomError(response.ErrBadRequest, "Status must be 0 or 1", 400)
	}

	typeProduct.Status = req.Status

	err = t.repo.UpdateStatus(ctx, id, req.Status)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "Failed to update status type product", 500)
	}

	t.InvalidateTypeProductCache(ctx)

	return nil
}

// GetAllTypeProductCustomer implements services.TypeProductService.
func (t *TypeProductService) GetAllTypeProductCustomer(ctx context.Context) ([]*entities.TypeProductEntity, error) {
	redisKey := "type_products_customer"

	cached, err := configs.GetRedis(ctx, redisKey)
	if err == nil && cached != "" {
		var cachedData struct {
			TypeProducts []*entities.TypeProductEntity `json:"type_products"`
			Total        int                           `json:"total"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedData); err != nil {
			return cachedData.TypeProducts, nil
		}
	}

	typeProducts, err := t.repo.GetAllTypeProductCustomer(ctx)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to get type product", 500)
	}

	dataToCache, _ := json.Marshal(map[string]interface{}{
		"type_products": typeProducts,
	})

	_ = configs.SetRedis(ctx, redisKey, dataToCache, time.Minute*30)

	return typeProducts, nil
}

// GetByIdTypeProductCustomer implements services.TypeProductService.
func (t *TypeProductService) GetByIdTypeProductCustomer(ctx context.Context, id uuid.UUID) (*entities.TypeProductEntity, error) {
	redisKey := fmt.Sprintf("type_product_customer:%s", id.String())

	cached, err := configs.GetRedis(ctx, redisKey)
	if err == nil && cached != "" {
		var typeProduct entities.TypeProductEntity
		if err := json.Unmarshal([]byte(cached), &typeProduct); err != nil {
			return &typeProduct, nil
		}
	}

	typeProduct, err := t.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "Not Found Type Product", 404)
		}
		return nil, err
	}
	dataToCache, _ := json.Marshal(typeProduct)
	_ = configs.SetRedis(ctx, redisKey, dataToCache, time.Minute*30)

	return typeProduct, nil
}

func (b *TypeProductService) InvalidateTypeProductCache(ctx context.Context) {
	iter := b.rdb.Scan(ctx, 0, "type_products:*", 0).Iterator()
	for iter.Next(ctx) {
		b.rdb.Del(ctx, iter.Val())
	}

	iterID := b.rdb.Scan(ctx, 0, "type_product:*", 0).Iterator()
	for iterID.Next(ctx) {
		b.rdb.Del(ctx, iterID.Val())
	}
}
