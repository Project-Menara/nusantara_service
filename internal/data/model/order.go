package model

import "github.com/google/uuid"

type OrderType string

const (
	TakeAway OrderType = "TakeAway"
	Delivery OrderType = "Delivery"
)

type Order struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
}
