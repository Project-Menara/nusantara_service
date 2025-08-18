package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name          string                  `form:"name" json:"name"`
	Code          string                  `form:"code" json:"code"`
	Price         int                     `form:"price" json:"price"`
	Unit          string                  `form:"unit" json:"unit"`
	Description   string                  `form:"description" json:"description"`
	Status        int                     `form:"status" json:"status"`
	TypeProductID uuid.UUID               `form:"type_product_id" json:"type_product_id"`
	Cover         *multipart.FileHeader   `form:"cover" swaggerignore:"true"`
	Gallery       []*multipart.FileHeader `form:"gallery" swaggerignore:"true"`
}

type UpdateProductRequest struct {
	ID             uuid.UUID               `json:"-"`
	Name           *string                 `form:"name" json:"name"`
	Code           *string                 `form:"code" json:"code"`
	Price          *int                    `form:"price" json:"price"`
	Unit           *string                 `form:"unit" json:"unit"`
	Description    *string                 `form:"description" json:"description"`
	TypeProductID  *uuid.UUID              `form:"type_product_id" json:"type_product_id"`
	NewCover       *multipart.FileHeader   `form:"new_cover" swaggerignore:"true"`
	NewGallery     []*multipart.FileHeader `form:"new_gallery" swaggerignore:"true"`
	ReplaceGallery bool                    `form:"replace_gallery" json:"replace_gallery"`
}

type UpdatStatusProducRequest struct {
	Status int `json:"status"`
}
