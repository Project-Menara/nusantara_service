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

type BannerService struct {
	repo       repositories.BannerRepository
	rdb        *redis.Client
	cloudinary cloudinary.CloudinaryService
}

func NewBannerUsecase(repo repositories.BannerRepository, rdb *redis.Client, cloudinary *cloudinary.CloudinaryService) services.BannerService {
	return &BannerService{repo: repo, rdb: rdb, cloudinary: *cloudinary}
}

// CreateBanner implements services.BannerService.
func (b *BannerService) CreateBanner(ctx context.Context, userId string, req dto.CreateBannerRequest, image *multipart.FileHeader) (*entities.BannerEntity, error) {
	existing, err := b.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)

	if existing != nil {
		return nil, response.NewCustomError(response.ErrExists, "name already exists", 409)
	}

	if name == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "name is required", 400)
	}

	if description == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "description is required", 400)
	}

	if image == nil {
		return nil, response.NewCustomError(response.ErrBadRequest, "image is required", 400)
	}

	if req.Status < 0 || req.Status > 1 {
		return nil, response.NewCustomError(response.ErrBadRequest, "Status must be 0 & 1", 400)
	}

	user, err := b.repo.FindByUserIDSuperAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	folder := "nusantara_service/banner"
	filename := fmt.Sprintf("banner_%s", req.Name) // id UUID unik banner

	imageUrl, err := b.cloudinary.UploadImage(ctx, image, folder, filename)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to upload image", 500)
	}

	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, response.NewCustomError(response.ErrBadRequest, "invalid user ID", 400)
	}

	newBanner := &entities.BannerEntity{
		Name:        name,
		Description: description,
		Status:      req.Status,
		Photo:       imageUrl,
		UserID:      userUUID,
		User:        *user,
	}

	_, err = b.repo.Create(ctx, newBanner)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to create banner", 500)
	}

	b.InvalidateBannerCache(ctx)

	return newBanner, nil

}

// GetAllBanner implements services.BannerService.
func (b *BannerService) GetAllBanner(ctx context.Context, page int, limit int, search string) ([]*entities.BannerEntity, int, error) {
	cacheKey := fmt.Sprintf("banner:search:%s:page:%d:limit:%d", search, page, limit)
	cached, err := configs.GetRedis(ctx, cacheKey)
	if err == nil {
		var result struct {
			Data  []*entities.BannerEntity `json:"data"`
			Total int                      `json:"total"`
		}

		_ = json.Unmarshal([]byte(cached), &result)
		return result.Data, result.Total, nil
	}

	offset := (page - 1) * limit
	// redisKey := fmt.Sprintf("banners:page=%d&limit=%d", page, limit)

	// Cek apakah data ada di redis
	// cached, err := configs.GetRedis(ctx, redisKey)
	// if err == nil && cached != "" {
	// 	var cachedData struct {
	// 		Banners []*entities.BannerEntity `json:"banners"`
	// 		Total   int                      `json:"total"`
	// 	}
	// 	if err := json.Unmarshal([]byte(cached), &cachedData); err == nil {
	// 		return cachedData.Banners, cachedData.Total, nil
	// 	}
	// }

	banners, err := b.repo.FindAllWithSearch(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrNotFound, "failed to find banner", 404)
	}

	total, err := b.repo.CountAllWithSearch(ctx, search)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrBadRequest, "failed to count banner", 400)
	}
	dataToCache, _ := json.Marshal(map[string]interface{}{
		"data":  banners,
		"total": total,
	})

	_ = configs.SetRedis(ctx, cacheKey, dataToCache, time.Minute*10)

	return banners, total, nil
}

// GetByIdBanner implements services.BannerService.
func (b *BannerService) GetByIdBanner(ctx context.Context, id uuid.UUID) (*entities.BannerEntity, error) {
	redisKey := fmt.Sprintf("banner:%s", id.String())

	// Cek dari redis
	cached, err := configs.GetRedis(ctx, redisKey)
	if err == nil && cached != "" {
		var banner entities.BannerEntity
		if err := json.Unmarshal([]byte(cached), &banner); err == nil {
			return &banner, nil
		}
	}

	banner, err := b.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "banner not found", 404)
		}
		return nil, err
	}

	dataToCache, _ := json.Marshal(banner)
	_ = configs.SetRedis(ctx, redisKey, dataToCache, time.Minute*30)

	return banner, nil
}

// UpdateBanner implements services.BannerService.
func (b *BannerService) UpdateBanner(ctx context.Context, userId string, id uuid.UUID, req dto.UpdateBannerRequest, image *multipart.FileHeader) (*entities.BannerEntity, error) {
	banner, err := b.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "banner not found", 404)
		}
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)

	if name != "" {
		existing, _ := b.repo.FindByName(ctx, req.Name)
		if existing != nil && existing.ID != id {
			return nil, response.NewCustomError(response.ErrExists, "banner already exists", 409)
		}
		banner.Name = name
	}
	if description != "" {
		banner.Description = description
	}

	if image != nil {
		if banner.Photo != "" {
			publicID := utils.ExtractPublicIDFromCloudinaryURL(banner.Photo)
			if publicID != "" {
				if err := b.cloudinary.DestroyImage(ctx, publicID); err != nil {
					return nil, response.NewCustomError(response.ErrInternal, "failed to delete image", 500)
				}
			}
		}

		folder := "nusantara_service/banner"
		filename := fmt.Sprintf("banner_%s", id.String()) // id UUID unik banner

		imageUrl, err := b.cloudinary.UploadImage(ctx, image, folder, filename)

		if err != nil {
			return nil, response.NewCustomError(response.ErrInternal, "failed to upload image", 500)
		}

		banner.Photo = imageUrl
	}

	user, err := b.repo.FindByUserIDSuperAdmin(ctx, userId)
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

	banner.UserID = userUUID

	updatedBanner, err := b.repo.Update(ctx, id, banner)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "Failed to update banner", 500)
	}
	b.InvalidateBannerCache(ctx)

	return updatedBanner, nil

}

// DeleteBanner implements services.BannerService.
func (b *BannerService) DeleteBanner(ctx context.Context, id uuid.UUID) error {
	banner, err := b.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "banner not found", 404)
		}
		return err
	}

	if banner.Photo != "" {
		publicID := utils.ExtractPublicIDFromCloudinaryURL(banner.Photo)
		if publicID != "" {
			if err := b.cloudinary.DestroyImage(ctx, publicID); err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to delete image", 500)
			}
		}
	}
	b.InvalidateBannerCache(ctx)

	return b.repo.Delete(ctx, id)
}

func (b *BannerService) InvalidateBannerCache(ctx context.Context) {
	iter := b.rdb.Scan(ctx, 0, "banners:*", 0).Iterator()
	for iter.Next(ctx) {
		b.rdb.Del(ctx, iter.Val())
	}

	iterID := b.rdb.Scan(ctx, 0, "banner:*", 0).Iterator()
	for iterID.Next(ctx) {
		b.rdb.Del(ctx, iterID.Val())
	}
}

// UpdateStatusBanner implements services.BannerService.
func (b *BannerService) UpdateStatusBanner(ctx context.Context, bannerId string, req dto.UpdateStatusBannerRequest) error {
	uuID, err := uuid.Parse(bannerId)
	if err != nil {
		return err
	}
	banner, err := b.repo.FindById(ctx, uuID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "banner not found", 404)
		}
		return err
	}
	if req.Status < 0 || req.Status > 1 {
		return response.NewCustomError(response.ErrBadRequest, "Status must be 0 or 1", 400)
	}

	banner.Status = req.Status

	err = b.repo.UpdateStatus(ctx, bannerId, req.Status)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "Failed to update status banner", 500)
	}

	b.InvalidateBannerCache(ctx)

	return nil
}

// GetAllBannerCustomer implements services.BannerService.
func (b *BannerService) GetAllBannerCustomer(ctx context.Context) ([]*entities.BannerEntity, error) {
	banners, err := b.repo.GetAllBannerCustomer(ctx)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "Failed to get banner", 500)
	}

	return banners, nil
}

// GetByIdBannerCustomer implements services.BannerService.
func (b *BannerService) GetByIdBannerCustomer(ctx context.Context, id uuid.UUID) (*entities.BannerEntity, error) {
	redisKey := fmt.Sprintf("banner_customer:%s", id.String())

	cached, err := configs.GetRedis(ctx, redisKey)
	if err == nil && cached != "" {
		var banner entities.BannerEntity
		if err := json.Unmarshal([]byte(cached), &banner); err != nil {
			return &banner, nil
		}
	}

	banner, err := b.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "Not Found banner", 404)
		}
		return nil, err
	}

	dataToCache, _ := json.Marshal(banner)
	_ = configs.SetRedis(ctx, redisKey, dataToCache, time.Minute*30)

	return banner, nil
}
