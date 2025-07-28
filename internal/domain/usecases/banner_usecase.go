package usecases

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	"strings"

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

	folder := fmt.Sprintf("nusantara_service/banner/%s", req.Name)
	imageUrl, err := b.cloudinary.UploadImage(ctx, image, folder)
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

	return newBanner, nil

}

// GetAllBanner implements services.BannerService.
func (b *BannerService) GetAllBanner(ctx context.Context, page int, limit int) ([]*entities.BannerEntity, int, error) {
	offset := (page - 1) * limit
	banners, err := b.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrNotFound, "failed to find banner", 404)
	}

	total, err := b.repo.CountAll(ctx)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrBadRequest, "failed to count banner", 400)
	}

	return banners, total, nil
}

// GetByIdBanner implements services.BannerService.
func (b *BannerService) GetByIdBanner(ctx context.Context, id uuid.UUID) (*entities.BannerEntity, error) {
	banner, err := b.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "banner not found", 404)
		}
		return nil, err
	}

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
			if publicID == "" {
				if err := b.cloudinary.DestroyImage(ctx, publicID); err != nil {
					return nil, response.NewCustomError(response.ErrInternal, "failed to delete image", 500)
				}
			}
		}
		folder := fmt.Sprintf("nusantara_service/banner/%s", req.Name)
		imageUrl, err := b.cloudinary.UploadImage(ctx, image, folder)
		if err != nil {
			return nil, response.NewCustomError(response.ErrInternal, "failed to upload image", 500)
		}

		banner.Photo = imageUrl
	}

	_, err = b.repo.FindByUserIDSuperAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	updatedBanner, err := b.repo.Update(ctx, id, banner)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "Failed to update banner", 500)
	}

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
		if publicID == "" {
			if err := b.cloudinary.DestroyImage(ctx, publicID); err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to delete image", 500)
			}
		}
	}

	return b.repo.Delete(ctx, id)
}
