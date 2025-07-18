package entities

import "github.com/google/uuid"

type RoleEntity struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"type:varchar(100);unique;not null"`
}

func (RoleEntity) TableName() string {
	return "roles"
}
