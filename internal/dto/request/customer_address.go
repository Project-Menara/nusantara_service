package dto

import "github.com/google/uuid"

type CreateAddressRequest struct {
	Label       string   `json:"label" binding:"required"`
	AddressText string   `json:"address_text" binding:"required"`
	Lat         *float64 `json:"lat" binding:"required"`
	Lng         *float64 `json:"lang" binding:"required"`
}

type UpdateAddressRequest struct {
	Label       string   `json:"label"`
	AddressText string   `json:"address_text"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
}

type NearbyShopsRequest struct {
	Lat       float64    `json:"lat,omitempty"`
	Lng       float64    `json:"lng,omitempty"`
	AddressID *uuid.UUID `json:"address_id,omitempty"`
	RadiusKm  float64    `json:"radius_km"`
}
