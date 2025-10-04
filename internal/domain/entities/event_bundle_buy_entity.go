package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventBundleBuyEntity struct {
	ID        uuid.UUID      `json:"id"`
	EventID   uuid.UUID      `json:"-"`
	Event     EventEntity    `json:"event"`
	ProductID uuid.UUID      `json:"-"`
	Product   ProductEntity  `json:"product"`
	Quantity  int            `json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (EventBundleBuyEntity) TableName() string {
	return "event_bundle_buys"
}
