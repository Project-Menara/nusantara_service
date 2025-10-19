package cashierrepo

import (
	"context"
	"nusantara_service/internal/domain/entities"
	shopresponse "nusantara_service/internal/dto/responses/shop_response"

	"github.com/google/uuid"
)

type ICashierRepository interface {
	GetShop(ctx context.Context, cashierID uuid.UUID) ([]shopresponse.ShopNameCashierResponse, error)
	FindCashier(ctx context.Context, userID string) (*entities.UserEntity, error)
	GetDetailShop(ctx context.Context, cashierID uuid.UUID, shopID uuid.UUID) (*shopresponse.ShopCashierResponse, error)
	GetDetailProduct(ctx context.Context, cashierID uuid.UUID, shopID uuid.UUID) ([]shopresponse.ShopProductCashierResponse, error)
}
