package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerAddressEntity struct {
	ID          uuid.UUID      `json:"id"`
	UserID      uuid.UUID      `json:"-"`
	User        UserEntity     `json:"user"`
	Label       string         `json:"label"`
	AddressText string         `json:"address_text"`
	Lat         float64        `json:"lat"`
	Lng         float64        `json:"lang"`
	IsDefault   bool           `json:"is_default"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (CustomerAddressEntity) TableName() string {
	return "customer_addresses"
}
