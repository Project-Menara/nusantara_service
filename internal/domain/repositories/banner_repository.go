package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"

	"github.com/google/uuid"
)

type BannerRepository interface {
	Create(ctx context.Context, banner *entities.BannerEntity) (*entities.BannerEntity, error)
	FindAll(ctx context.Context, limit, offset int) ([]*entities.BannerEntity, error)
	FindById(ctx context.Context, id uuid.UUID) (*entities.BannerEntity, error)
	Update(ctx context.Context, id uuid.UUID, data *entities.BannerEntity) (*entities.BannerEntity, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, bannerId string, status int) error
	GetAllBannerCustomer(ctx context.Context) ([]*entities.BannerEntity, error)

	FindByName(ctx context.Context, name string) (*entities.BannerEntity, error)
	CountAll(ctx context.Context) (int, error)
	FindByUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error)

	FindAllWithSearch(ctx context.Context, limit, offset int, search string) ([]*entities.BannerEntity, error)
	CountAllWithSearch(ctx context.Context, search string) (int, error)
}
