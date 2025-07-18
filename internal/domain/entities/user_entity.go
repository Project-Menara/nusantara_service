package entities

import "time"

type UserEntity struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"type:varchar(255)"`
	Username  string    `gorm:"type:varchar(255);unique"`
	Email     string    `gorm:"type:varchar(255);unique"`
	Phone     string    `gorm:"type:varchar(255);unique"`
	Password  string    `gorm:"type:varchar(255)"`
	Gender    string    `gorm:"type:varchar(100)"`
	Birth     time.Time `gorm:"type:date"`
	Photo     string    `gorm:"type:varchar(100)"`
	RoleID    string    `gorm:"type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
