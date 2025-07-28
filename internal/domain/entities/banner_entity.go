package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BannerEntity struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Photo       string     `gorm:"type:varchar(255)" json:"photo"`
	Name        string     `gorm:"type:varchar(255);unique;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Status      int        `gorm:"type:int;not null" json:"status"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"-"`
	User        UserEntity `gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (BannerEntity) TableName() string {
	return "banners"
}
