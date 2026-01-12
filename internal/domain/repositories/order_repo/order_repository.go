package orderrepo

import (
	"context"
	"nusantara_service/internal/data/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	WithTx(tx *gorm.DB) OrderRepository
	CreateOrder(ctx context.Context, data *model.Order) error
	CreateOrderItem(ctx context.Context, data []model.OrderItem) error
	CreateOrderEvent(ctx context.Context, data *model.OrderEvent) error
	CreateOrderReward(ctx context.Context, data *model.OrderReward) error
	CreateOrderVoucher(ctx context.Context, data *model.OrderVoucher) error

	GetCartUser(ctx context.Context, userId uuid.UUID) (*model.Cart, error)
	GetDetailCartUser(ctx context.Context, cartId uuid.UUID) ([]*model.CartItem, error)
	// DeleteCartItemsByIDs deletes cart items by their IDs.
	DeleteCartItemsByIDs(ctx context.Context, ids []uuid.UUID) error
}
