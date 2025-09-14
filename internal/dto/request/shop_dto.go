package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type CreateShopRequest struct {
	Name        string                  `form:"name" json:"name"`
	Description string                  `form:"description" json:"description"`
	FullAddress string                  `form:"full_address" json:"full_address"`
	Lat         *float64                `form:"lat" json:"lat"`
	Lng         *float64                `form:"lang" json:"lang"`
	Status      *int                    `form:"status" json:"status"`
	Cover       *multipart.FileHeader   `form:"cover" swaggerignore:"true"`
	Gallery     []*multipart.FileHeader `form:"gallery" swaggerignore:"true"`
	CashierIDs  []uuid.UUID             `form:"cashier_ids" json:"cashier_ids"`
	Products    []AssignProductRequest  `form:"products" json:"products"`
}

type UpdateShopRequest struct {
	Name           string                  `form:"name" json:"name"`
	Description    string                  `form:"description" json:"description"`
	FullAddress    string                  `form:"full_address" json:"full_address"`
	Lat            float64                 `form:"lat" json:"lat"`
	Lng            float64                 `form:"lang" json:"lang"`
	NewCover       *multipart.FileHeader   `form:"cover" swaggerignore:"true"`
	NewGallery     []*multipart.FileHeader `form:"gallery" swaggerignore:"true"`
	ReplaceGallery bool                    `form:"replace_gallery" json:"replace_gallery"`

	// optional: replace assignments when present
	CashierIDs []uuid.UUID            `form:"cashier_ids" json:"cashier_ids"`
	Products   []AssignProductRequest `form:"products" json:"products"`
}

// // Untuk assign kasir ke toko
// type AssignCashierRequest struct {
// 	ShopID    uuid.UUID `json:"shop_id"`
// 	CashierID uuid.UUID `json:"cashier_id"`
// }

// Untuk assign produk ke toko
type AssignProductRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Price     *float64  `json:"price,omitempty"`
	Stock     int       `json:"stock"`
	Status    *int      `json:"status,omitempty"`
}

type UpdateStatusShopRequest struct {
	Status int `json:"status"`
}
