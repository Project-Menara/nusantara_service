package dto

import "github.com/google/uuid"

type AddCartItemRequest struct {
	ProductID uuid.UUID `form:"product_id" json:"product_id"`
}

type UpdateSelectedCartItem struct {
	ProductID uuid.UUID `form:"product_id" json:"product_id"`
	Selected  bool      `form:"selected" json:"selected"`
}
