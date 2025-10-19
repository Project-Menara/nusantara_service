package cashierservice

import (
	"context"
	shopresponse "nusantara_service/internal/dto/responses/shop_response"

	"github.com/google/uuid"
)

type CashierService interface {
	GetAllNameShop(ctx context.Context, cashierID uuid.UUID) ([]shopresponse.ShopNameCashierResponse, error)
	GetDetailShopCashier(ctx context.Context, cashierID uuid.UUID, shopID uuid.UUID) (*shopresponse.ShopCashierResponse, error)
	GetDetailShopProduct(ctx context.Context, cashierID uuid.UUID, shopID uuid.UUID) ([]shopresponse.ShopProductCashierResponse, error)
}
