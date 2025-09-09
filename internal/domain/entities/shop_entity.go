package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopEntity struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Cover       string         `json:"cover"`
	Description string         `json:"description"`
	FullAddress string         `json:"full_address"`
	Lat         float64        `json:"lat"`
	Lng         float64        `json:"lang"`
	Status      int            `json:"status"`
	CreatedBy   uuid.UUID      `json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"update_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`

	Creator      UserEntity          `json:"created_by" gorm:"foreignKey:CreatedBy;references:ID"`
	ShopImages   []ShopImageEntity   `json:"shop_images" gorm:"foreignKey:ShopID;references:ID"`
	ShopProducts []ShopProductEntity `json:"shop_products" gorm:"foreignKey:ShopID;references:ID"`
	ShopCashiers []ShopCashierEntity `json:"shop_cashiers" gorm:"foreignKey:ShopID;references:ID"`
}

type ShopImageEntity struct {
	ID        uuid.UUID      `json:"id"`
	ShopID    uuid.UUID      `json:"shop_id"`
	ImageID   uuid.UUID      `json:"-"`
	Altext    string         `json:"alt_text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	Shop  ShopEntity  `json:"-" gorm:"foreignKey:ShopID;references:ID"`
	Image ImageEntity `json:"image" gorm:"foreignKey:ImageID;references:ID"`
}

func (ShopEntity) TableName() string {
	return "shops"
}

func (ShopImageEntity) TableName() string {
	return "shop_images"
}
