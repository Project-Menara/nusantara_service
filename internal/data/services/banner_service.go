package services

import (
	"context"
	"mime/multipart"
	"nusantara_service/internal/domain/entities"
	dto "nusantara_service/internal/dto/request"

	"github.com/google/uuid"
)

type BannerService interface {
	CreateBanner(ctx context.Context, userId string, req dto.CreateBannerRequest, image *multipart.FileHeader) (*entities.BannerEntity, error)
	GetAllBanner(ctx context.Context, page, limit int, search string) ([]*entities.BannerEntity, int, error)
	GetByIdBanner(ctx context.Context, id uuid.UUID) (*entities.BannerEntity, error)
	UpdateBanner(ctx context.Context, userId string, id uuid.UUID, req dto.UpdateBannerRequest, image *multipart.FileHeader) (*entities.BannerEntity, error)
	DeleteBanner(ctx context.Context, id uuid.UUID) error
	UpdateStatusBanner(ctx context.Context, bannerId string, req dto.UpdateStatusBannerRequest) error

	GetAllBannerCustomer(ctx context.Context) ([]*entities.BannerEntity, error)
	GetByIdBannerCustomer(ctx context.Context, id uuid.UUID) (*entities.BannerEntity, error)
}
