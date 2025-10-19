package cashierusecase

import (
	"context"
	"errors"
	cashierservice "nusantara_service/internal/data/services/CashierService"
	cashierrepo "nusantara_service/internal/domain/repositories/CashierRepo"
	shopresponse "nusantara_service/internal/dto/responses/shop_response"
	"nusantara_service/internal/response"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CashierService struct {
	cashierRepo cashierrepo.ICashierRepository
	rdb         *redis.Client
}

func NewCashierUsecase(cashierRepo cashierrepo.ICashierRepository, rdb *redis.Client) cashierservice.CashierService {
	return &CashierService{cashierRepo: cashierRepo, rdb: rdb}
}

// GetAllNameShop implements cashierservice.CashierService.
func (c *CashierService) GetAllNameShop(ctx context.Context, cashierID uuid.UUID) ([]shopresponse.ShopNameCashierResponse, error) {
	cashier, err := c.cashierRepo.FindCashier(ctx, cashierID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "cashier not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get cashier", 500)
	}
	if cashier.Role.Name != "admin" {
		return nil, response.NewCustomError(response.ErrForbidden, "not permission", 403)
	}
	shopNames, err := c.cashierRepo.GetShop(ctx, uuid.MustParse(cashier.ID))
	if err != nil {
		return nil, response.NewCustomError(response.ErrNotFound, "Shop Not Found", 404)
	}

	return shopNames, nil
}

// GetDetailShopCashier implements cashierservice.CashierService.
func (c *CashierService) GetDetailShopCashier(ctx context.Context, cashierID uuid.UUID, shopID uuid.UUID) (*shopresponse.ShopCashierResponse, error) {
	cashier, err := c.cashierRepo.FindCashier(ctx, cashierID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "cashier not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get cashier", 500)
	}
	if cashier.Role.Name != "admin" {
		return nil, response.NewCustomError(response.ErrForbidden, "not permission", 403)
	}

	shopCashier, err := c.cashierRepo.GetDetailShop(ctx, uuid.MustParse(cashier.ID), shopID)
	if err != nil {
		return nil, response.NewCustomError(response.ErrNotFound, "Shop Not Found", 404)
	}

	return shopCashier, nil

}

// GetDetailShopProduct implements cashierservice.CashierService.
func (c *CashierService) GetDetailShopProduct(ctx context.Context, cashierID uuid.UUID, shopID uuid.UUID) ([]shopresponse.ShopProductCashierResponse, error) {
	cashier, err := c.cashierRepo.FindCashier(ctx, cashierID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "cashier not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get cashier", 500)
	}
	if cashier.Role.Name != "admin" {
		return nil, response.NewCustomError(response.ErrForbidden, "not permission", 403)
	}

	shopProductCashier, err := c.cashierRepo.GetDetailProduct(ctx, uuid.MustParse(cashier.ID), shopID)
	if err != nil {
		return nil, response.NewCustomError(response.ErrNotFound, "Product Not Found", 404)
	}

	return shopProductCashier, nil
}
