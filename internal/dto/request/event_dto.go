package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type CreateEventRequest struct {
	Name               string                 `form:"name" json:"name"`
	Type               string                 `form:"type_event" json:"type_event"`
	StartDate          time.Time              `form:"start_date" json:"start_date"`
	EndDate            time.Time              `form:"end_date" json:"end_date"`
	Cover              *multipart.FileHeader  `form:"cover" swaggerignore:"true"`
	Status             int                    `form:"status" json:"status"`
	EventProducts      []AddEventProduct      `form:"event_products" json:"event_products"`
	EventBundleBuys    []AddEventBundleBuy    `form:"event_bundle_buys" json:"event_bundle_buys"`
	EventBundleRewards []AddEventBundleReward `form:"event_bundle_rewards" json:"event_bundle_rewards"`
}

type UpdateEventRequest struct {
	Name               string                 `form:"name" json:"name"`
	Type               string                 `form:"type_event" json:"type_event"`
	StartDate          time.Time              `form:"start_date" json:"start_date"`
	EndDate            time.Time              `form:"end_date" json:"end_date"`
	NewCover           *multipart.FileHeader  `form:"new_cover" swaggerignore:"true"`
	EventProducts      []AddEventProduct      `form:"event_products" json:"event_products"`
	EventBundleBuys    []AddEventBundleBuy    `form:"event_bundle_buys" json:"event_bundle_buys"`
	EventBundleRewards []AddEventBundleReward `form:"event_bundle_rewards" json:"event_bundle_rewards"`
}

type UpdateStatusEventRequest struct {
	Status int `json:"status"`
}

type AddEventProduct struct {
	ProductID       uuid.UUID `json:"product_id"`
	DiscountPercent int       `json:"discount_percent"`
}
type AddEventBundleBuy struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}
type AddEventBundleReward struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}
