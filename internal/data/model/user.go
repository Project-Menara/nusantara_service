package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name     string    `gorm:"type:varchar(255)"`
	Username string    `gorm:"type:varchar(255);unique;not null"`
	Email    string    `gorm:"type:varchar(255);unique;not null"`
	Phone    string    `gorm:"type:varchar(255);unique"`
	Password string    `gorm:"type:varchar(255);not null"`
	Gender   string    `gorm:"type:varchar(100);not null"`
	Birth    time.Time `gorm:"type:date;not null"`
	Photo    string    `gorm:"type:varchar(100)"`

	RoleID uuid.UUID `gorm:"type:uuid"`
	Role   Role      `gorm:"foreignKey:RoleID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
