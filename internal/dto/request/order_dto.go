package dto

import "github.com/google/uuid"

type CreateOrderRequest struct {
	CartID uuid.UUID                `json:"cart_id"`
	ShopID uuid.UUID                `json:"shop_id"`
	Items  []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type CreateOrderEvent struct {
}
