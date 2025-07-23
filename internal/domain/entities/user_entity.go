package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserEntity struct {
	ID          string         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255)" json:"name"`
	Username    string         `gorm:"type:varchar(255);unique" json:"username"`
	Email       string         `gorm:"type:varchar(255);unique" json:"email"`
	Phone       *string        `gorm:"type:varchar(255);unique" json:"phone"`
	Password    string         `gorm:"type:varchar(255)" json:"password"`
	Gender      *string        `gorm:"type:varchar(100)" json:"gender"`
	DateOfBirth *time.Time     `gorm:"type:date" json:"date_of_birth"`
	Photo       *string        `gorm:"type:varchar(100)" json:"photo"`
	RoleID      uuid.UUID      `gorm:"type:uuid" json:"-"`
	Role        RoleEntity     `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (UserEntity) TableName() string {
	return "users"
}
