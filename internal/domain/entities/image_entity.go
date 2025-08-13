package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageEntity struct {
	ID        uuid.UUID      `json:"id"`
	ImagePath string         `json:"image_path"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (ImageEntity) TableName() string {
	return "images"
}
