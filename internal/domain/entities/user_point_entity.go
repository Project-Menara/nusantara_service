package entities

import (
	"nusantara_service/internal/data/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPointEntity struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;uniqueIndex" json:"-"`
	User        model.User     `gorm:"foreignKey:UserID" json:"user"`
	TotalPoints int            `gorm:"type:int" json:"total_points"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (UserPointEntity) TableName() string {
	return "user_points"
}
