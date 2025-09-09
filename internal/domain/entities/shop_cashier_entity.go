package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopCashierEntity struct {
	ID         uuid.UUID `json:"id"`
	ShopID     uuid.UUID `json:"-"`
	CashierID  uuid.UUID `json:""`
	AssignedAt time.Time `json:"assign_at"`

	Shop ShopEntity `gorm:"foreignKey:ShopID;references:ID"`
	User UserEntity `gorm:"foreignKey:CashierID;references:ID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (ShopCashierEntity) TableName() string {
	return "shop_cashiers"
}
