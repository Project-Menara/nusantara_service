package entities

import "github.com/google/uuid"

type RoleEntity struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(100);unique;not null" json:"name"`
}

func (RoleEntity) TableName() string {
	return "roles"
}
