package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductEntity struct {
	ID            uuid.UUID            `json:"id"`
	Name          string               `json:"name"`
	ImageID       uuid.UUID            `json:"-"`
	Image         ImageEntity          `json:"image" gorm:"foreignKey:ImageID;references:ID"`
	Code          string               `json:"code"`
	Price         int                  `json:"price"`
	Unit          string               `json:"unit"`
	Description   string               `json:"description"`
	Status        int                  `json:"status"`
	TypeProductID uuid.UUID            `json:"-"`
	TypeProduct   TypeProductEntity    `json:"type_product" gorm:"foreignKey:TypeProductID;references:ID"`
	ProductImages []ProductImageEntity `json:"product_images" gorm:"foreignKey:ProductID;references:ID"`
	CreatedBy     uuid.UUID            `json:"-"`
	User          UserEntity           `json:"created_by" gorm:"foreignKey:CreatedBy;references:ID"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	DeletedAt     gorm.DeletedAt       `json:"deleted_at"`
}

type ProductImageEntity struct {
	ID        uuid.UUID      `json:"id"`
	ProductID uuid.UUID      `json:"product_id"`
	Product   ProductEntity  `json:"-" gorm:"foreignKey:ProductID;references:ID"`
	ImageID   uuid.UUID      `json:""`
	Image     ImageEntity    `json:"image" gorm:"foreignKey:ImageID;references:ID"`
	AltText   string         `json:"alt_text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (ProductEntity) TableName() string {
	return "products"
}

func (ProductImageEntity) TableName() string {
	return "product_images"
}
