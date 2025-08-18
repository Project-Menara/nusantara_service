package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entities.ProductEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ProductEntity, error)
	GetAll(ctx context.Context, offset, limit int, search string) ([]*entities.ProductEntity, int, error)
	Update(ctx context.Context, id uuid.UUID, product *entities.ProductEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, productID uuid.UUID, status int) error

	FindByUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error)
	FindByName(ctx context.Context, name string) (*entities.ProductEntity, error)

	//images
	CreateImages(ctx context.Context, images []entities.ImageEntity) error
	CreateProductImages(ctx context.Context, piv []entities.ProductImageEntity) error
	DeleteProductImages(ctx context.Context, productID uuid.UUID) error
	GetProductImagePublicIDs(ctx context.Context, productID uuid.UUID) ([]string, error)
}
