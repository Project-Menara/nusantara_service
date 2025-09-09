package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"

	"github.com/google/uuid"
)

type ShopRepository interface {
	Create(ctx context.Context, data *entities.ShopEntity) error
	FindAll(ctx context.Context, offset, limit int, search string) ([]*entities.ShopEntity, int, error)
	FindById(ctx context.Context, id uuid.UUID) (*entities.ShopEntity, error)
	Update(ctx context.Context, id uuid.UUID, data *entities.ShopEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status int) error

	FindUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error)
	FindByName(ctx context.Context, name string) (*entities.ShopEntity, error)
	FindCashier(ctx context.Context, cashierID uuid.UUID) (*entities.UserEntity, error)

	CreateImage(ctx context.Context, img *entities.ImageEntity) error
	UpdateShopCover(ctx context.Context, shopID uuid.UUID, coverURL string) error
	CreateGallery(ctx context.Context, shopID uuid.UUID, imageID uuid.UUID, alt string) error
	DeleteGalleryByShopID(ctx context.Context, shopID uuid.UUID) error
	AssignCashiers(ctx context.Context, shopID uuid.UUID, cashierIDs []uuid.UUID) error
	AssignProducts(ctx context.Context, shopID uuid.UUID, items []entities.ShopProductEntity) error
	DeleteShopCashiersByShopID(ctx context.Context, shopID uuid.UUID) error
	DeleteShopProductsByShopID(ctx context.Context, shopID uuid.UUID) error
	ReplaceCashiers(ctx context.Context, shopID uuid.UUID, cashierIDs []uuid.UUID) error
	ReplaceProducts(ctx context.Context, shopID uuid.UUID, items []entities.ShopProductEntity) error
}
