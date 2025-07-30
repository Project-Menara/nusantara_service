package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TypeProductEntity struct {
	ID     uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Image  string     `gorm:"type:varchar(255)" json:"image"`
	Name   string     `gorm:"type:varchar(255);not null" json:"name"`
	Status int        `gorm:"type:int;not null" json:"status"`
	UserID uuid.UUID  `gorm:"type:uuid" json:"-"`
	User   UserEntity `gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (TypeProductEntity) TableName() string {
	return "type_products"
}
