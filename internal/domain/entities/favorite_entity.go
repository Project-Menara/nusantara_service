package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FavoriteEntity struct {
	ID     uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID uuid.UUID  `gorm:"type:uuid" json:"-"`
	User   UserEntity `gorm:"foreignKey:UserID;OnDelete:CASCADE" json:"user"`

	CreatedAt     time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time            `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt       `gorm:"index" json:"deleted_at"`
	FavoriteItems []FavoriteItemEntity `gorm:"foreignKey:FavoriteID;references:ID"`
}

func (FavoriteEntity) TableName() string {
	return "favorites"
}
