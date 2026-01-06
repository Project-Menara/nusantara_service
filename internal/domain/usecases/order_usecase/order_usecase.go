package orderusecase

import (
	"context"
	"errors"
	"nusantara_service/internal/data/model"
	orderrepoimpl "nusantara_service/internal/data/repositories/order_repo_impl"
	orderservice "nusantara_service/internal/data/services/order_service"
	"nusantara_service/internal/domain/repositories"
	orderrepo "nusantara_service/internal/domain/repositories/order_repo"
	"nusantara_service/internal/response"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrderUsecase struct {
	orderRepo   orderrepo.OrderRepository
	userRepo    repositories.UserRepository
	productRepo repositories.ProductRepository
	rdb         *redis.Client
}

func NewOrderUsecase(orderRepo orderrepo.OrderRepository, userRepo repositories.UserRepository, productRepo repositories.ProductRepository, rdb *redis.Client) orderservice.OrderService {
	return &OrderUsecase{orderRepo: orderRepo, userRepo: userRepo, productRepo: productRepo, rdb: rdb}
}

func generateCodeOrder() string {
	id := strings.ToUpper(uuid.New().String())
	return "ORD-" + id[:8]
}

// CreateOrder implements [orderservice.OrderService].
func (o *OrderUsecase) CreateOrder(ctx context.Context, customerID uuid.UUID) error {
	// 1. Get Customer
	customer, err := o.userRepo.FindUserById(ctx, customerID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "customer not found", 404)
		}
		return response.NewCustomError(response.ErrInternal, "failed to get customer", 500)
	}

	// 2. Get Cart
	cart, err := o.orderRepo.GetCartUser(ctx, uuid.MustParse(customer.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "cart not found", 404)
		}
		return response.NewCustomError(response.ErrInternal, "failed to get cart", 500)
	}

	// 3. Persiapkan Data Order Utama
	newOrder := &model.Order{
		ID:       uuid.New(),
		Code:     generateCodeOrder(),
		UserID:   uuid.MustParse(customer.ID),
		ShopID:   cart.ShopID,
		CartID:   cart.ID,
		Status:   model.OrderDraft,
		SubTotal: 0, // Akan diupdate setelah loop item
		Total:    0,
	}

	// 4. Jalankan Transaksi
	orderImpl := o.orderRepo.(*orderrepoimpl.OrderRepositoryImpl)
	return orderImpl.DB().Transaction(func(tx *gorm.DB) error {
		// PENTING: Gunakan repo dengan instance transaksi 'tx'
		txRepo := o.orderRepo.WithTx(tx)

		// Simpan Header Order
		if err := txRepo.CreateOrder(ctx, newOrder); err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to create order", 500)
		}

		cartItems, err := txRepo.GetDetailCartUser(ctx, cart.ID)
		if err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to get cart items", 500)
		}

		var orderItems []model.OrderItem
		var totalOrderPrice float64

		for _, cp := range cartItems {
			if !cp.Selected {
				continue
			}

			// Idealnya: Gunakan Batch Fetch Product di luar loop untuk performa lebih baik
			cartProduct, err := o.productRepo.GetByID(ctx, cp.ProductID)
			if err != nil {
				return response.NewCustomError(response.ErrNotFound, "product not found", 404)
			}

			quantity := 1
			// Gunakan Quantity dari Cart, jangan hardcode 1
			itemSubTotal := float64(quantity) * float64(cartProduct.Price)
			totalOrderPrice += itemSubTotal

			orderItems = append(orderItems, model.OrderItem{
				ID:        uuid.New(),
				OrderID:   newOrder.ID,
				ProductID: cp.ProductID,
				Quantity:  quantity,
				SubTotal:  itemSubTotal,
			})
		}

		if len(orderItems) == 0 {
			return response.NewCustomError(response.ErrBadRequest, "no items selected", 400)
		}

		// Simpan Bulk Order Items
		if err := txRepo.CreateOrderItem(ctx, orderItems); err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to create order items", 500)
		}

		// Update total harga di tabel Order (Header)
		if err := tx.Model(newOrder).Updates(map[string]interface{}{
			"sub_total": totalOrderPrice,
			"total":     totalOrderPrice,
		}).Error; err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to update order total", 500)
		}

		return nil
	})
}
