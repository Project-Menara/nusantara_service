package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPointHistories struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID"`
	PointType   string         `gorm:"type:varchar(50);check:point_type IN ('reward','purchase','exchange')" json:"point_type"` // reward, purchase, exchange
	Source      string         `gorm:"type:varchar(255)" json:"source"`
	SourceId    string         `gorm:"type:varchar(255)" json:"source_id"`
	Points      int            `gorm:"type:int" json:"points"`
	ExpiredAt   *time.Time     `gorm:"type:date;index" json:"expired_at"`
	Description string         `gorm:"type:text" json:"description"`
	Direction   string         `gorm:"type:varchar(100);check:direction IN ('in','out')" json:"direction"` // in / out
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
