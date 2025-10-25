package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FavoriteItemEntity struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FavoriteID uuid.UUID `gorm:"type:uuid" json:"-"`
	ProductID  uuid.UUID `gorm:"type:uuid" json:"-"`
	Selected   bool      `gorm:"type:boolean" json:"selected"`

	Favorite FavoriteEntity `gorm:"foreignKey:FavoriteID;OnDelete:CASCADE"`
	Product  ProductEntity  `gorm:"foreignKey:ProductID;OnDelete:CASCADE"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (FavoriteItemEntity) TableName() string {
	return "favorite_items"
}
