package dto

import "github.com/google/uuid"

type AddFavoriteItemRequest struct {
	ProductID uuid.UUID `form:"product_id" json:"product_id"`
}
