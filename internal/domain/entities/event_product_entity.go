package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventProductEntity struct {
	ID              uuid.UUID      `json:"id"`
	EventID         uuid.UUID      `json:"-"`
	Event           EventEntity    `json:"event"`
	ProductID       uuid.UUID      `json:"-"`
	Product         ProductEntity  `json:"product"`
	DiscountPercent int            `json:"discount_percent"`
	DiscountAmount  float64        `json:"discount_amount"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"`
}

func (EventProductEntity) TableName() string {
	return "event_products"
}
