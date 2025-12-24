package cartresponse

import (
	"nusantara_service/internal/domain/entities"
	productresponse "nusantara_service/internal/dto/responses/product_response"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartResponse struct {
	ID        uuid.UUID          `json:"id"`
	User      string             `json:"user"`
	Shop      string             `json:"shop"`
	Status    int                `json:"status"`
	CartItems []CartItemResponse `json:"cart_items"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `json:"deleted_at"`
}

type CartItemResponse struct {
	ID       uuid.UUID                       `json:"id"`
	Selected bool                            `json:"selected"`
	Product  productresponse.ProductResponse `gorm:"foreignKey:ProductID;OnDelete:CASCADE"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func ToCartItemResponse(cartItem entities.CartItemEntity) CartItemResponse {
	productItem := productresponse.ToProductResponse(&cartItem.Product)
	return CartItemResponse{
		ID:        cartItem.ID,
		Selected:  cartItem.Selected,
		Product:   productItem,
		CreatedAt: cartItem.CreatedAt,
		UpdatedAt: cartItem.UpdatedAt,
		DeletedAt: cartItem.DeletedAt,
	}
}

func ToCartResponse(cart entities.CartEntity) CartResponse {
	cartItems := []CartItemResponse{}
	for _, cartItem := range cart.CartItems {
		cartItemRes := ToCartItemResponse(cartItem)
		cartItems = append(cartItems, cartItemRes)
	}
	return CartResponse{
		ID:        cart.ID,
		User:      cart.User.Name,
		Shop:      cart.Shop.Name,
		Status:    cart.Status,
		CartItems: cartItems,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
		DeletedAt: cart.DeletedAt,
	}
}
