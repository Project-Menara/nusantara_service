package orderrepoimpl

import (
	"context"
	"nusantara_service/internal/data/model"
	orderrepo "nusantara_service/internal/domain/repositories/order_repo"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func (p *OrderRepositoryImpl) DB() *gorm.DB {
	return p.db
}

func NewOrderRepositoryImpl(db *gorm.DB) orderrepo.OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

func (o *OrderRepositoryImpl) WithTx(tx *gorm.DB) orderrepo.OrderRepository {
	return &OrderRepositoryImpl{db: tx} // Mengembalikan repo baru dengan koneksi transaksi
}

// CreateOrder implements [orderrepo.OrderRepository].
func (o *OrderRepositoryImpl) CreateOrder(ctx context.Context, data *model.Order) error {
	return o.db.WithContext(ctx).Create(data).Error
}

// CreateOrderEvemt implements [orderrepo.OrderRepository].
func (o *OrderRepositoryImpl) CreateOrderEvent(ctx context.Context, data *model.OrderEvent) error {
	return o.db.WithContext(ctx).Create(data).Error
}

// CreateOrderItem implements [orderrepo.OrderRepository].
func (o *OrderRepositoryImpl) CreateOrderItem(ctx context.Context, data []model.OrderItem) error {
	return o.db.WithContext(ctx).Create(&data).Error
}

// CreateOrderReward implements [orderrepo.OrderRepository].
func (o *OrderRepositoryImpl) CreateOrderReward(ctx context.Context, data *model.OrderReward) error {
	return o.db.WithContext(ctx).Create(data).Error

}

// CreateOrderVoucher implements [orderrepo.OrderRepository].
func (o *OrderRepositoryImpl) CreateOrderVoucher(ctx context.Context, data *model.OrderVoucher) error {
	return o.db.WithContext(ctx).Create(data).Error

}

// GetCartUser implements [orderrepo.OrderRepository].
func (o *OrderRepositoryImpl) GetCartUser(ctx context.Context, userId uuid.UUID) (*model.Cart, error) {
	var cart model.Cart
	if err := o.db.WithContext(ctx).Where("user_id = ?", userId).First(&cart).Error; err != nil {
		return nil, err
	}

	return &cart, nil
}

// GetDetailCartUser implements [orderrepo.OrderRepository].
func (o *OrderRepositoryImpl) GetDetailCartUser(ctx context.Context, cartId uuid.UUID) ([]*model.CartItem, error) {
	var cartItems []*model.CartItem
	if err := o.db.WithContext(ctx).Where("cart_id = ?", cartId).Find(&cartItems).Error; err != nil {
		return nil, err
	}

	return cartItems, nil
}
