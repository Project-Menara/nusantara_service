package productresponse

import (
	"nusantara_service/internal/domain/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductResponse struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Image         string         `json:"image"`
	Code          string         `json:"code"`
	Price         int            `json:"price"`
	Unit          string         `json:"unit"`
	Description   string         `json:"description"`
	Status        int            `json:"status"`
	TypeProduct   string         `json:"type_product"`
	ProductImages []string       `json:"product_images"`
	CreatedBy     string         `json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

func ToProductResponse(product *entities.ProductEntity) ProductResponse {
	productImages := []string{}
	for _, productImage := range product.ProductImages {
		productImages = append(productImages, productImage.Image.ImagePath)
	}
	return ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Image:         product.Image.ImagePath,
		Code:          product.Code,
		Price:         product.Price,
		Unit:          product.Unit,
		Description:   product.Description,
		Status:        product.Status,
		TypeProduct:   product.TypeProduct.Name,
		ProductImages: productImages,
		CreatedBy:     product.User.Name,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
		DeletedAt:     product.DeletedAt,
	}
}
