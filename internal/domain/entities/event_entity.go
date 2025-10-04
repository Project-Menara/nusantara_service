package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventType string

const (
	// EventType values for TypeEvent
	EventTypeBundle EventType = "BUNDLE"
	EventTypeDiskon EventType = "DISKON"
)

type EventEntity struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name"`
	TypeEvent EventType `json:"type_event"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Cover     string    `json:"cover"`
	Status    int       `json:"status"`
	CreatedBy uuid.UUID `json:"created_by" gorm:"type:uuid"`

	User UserEntity `json:"user" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	EventProducts      []EventProductEntity      `json:"event_products" gorm:"foreignKey:EventID;references:ID"`
	EventBundleBuys    []EventBundleBuyEntity    `json:"event_bundle_buys" gorm:"foreignKey:EventID;references:ID"`
	EventBundleRewards []EventBundleRewardEntity `json:"event_bundle_rewards" gorm:"foreignKey:EventID;references:ID"`
}

func (EventEntity) TableName() string {
	return "events"
}
