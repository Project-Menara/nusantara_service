package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartEntity struct {
	ID        uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID        `gorm:"type:uuid" json:"-"`
	User      UserEntity       `gorm:"foreignKey:UserID;OnDelete:CASCADE" json:"user"`
	Status    int              `gorm:"type:int" json:"status"`
	ShopID    uuid.UUID        `gorm:"type:uuid" json:"-"`
	Shop      ShopEntity       `gorm:"foreignKey:ShopID;constraint:OnDelete:CASCADE;" json:"shop"`
	CreatedAt time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"deleted_at"`
	CartItems []CartItemEntity `gorm:"foreignKey:CartID;references:ID"`
}

func (CartEntity) TableName() string {
	return "carts"
}
