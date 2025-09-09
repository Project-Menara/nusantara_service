package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopProductEntity struct {
	ID         uuid.UUID      `json:"id"`
	ShopID     uuid.UUID      `json:"-"`
	ProductID  uuid.UUID      `json:"-"`
	Price      float64        `json:"price"`
	Stock      int            `json:"stock"`
	Status     int            `json:"status"`
	AssignedAt time.Time      `json:"assign_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`

	Shop    ShopEntity    `gorm:"foreignKey:ShopID;references:ID"`
	Product ProductEntity `gorm:"foreignKey:ProductID;references:ID"`
}

func (ShopProductEntity) TableName() string {
	return "shop_products"
}
