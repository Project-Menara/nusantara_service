package model

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"type:varchar(100);unique;not null"`
}
