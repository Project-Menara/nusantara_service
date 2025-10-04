package eventresponse

import (
	"nusantara_service/internal/domain/entities"
	productresponse "nusantara_service/internal/dto/responses/product_response"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventProductResponse struct {
	ID              uuid.UUID                       `json:"id"`
	Event           string                          `json:"event"`
	Product         productresponse.ProductResponse `json:"product"`
	DiscountPercent int                             `json:"discount_percent"`
	DiscountAmount  float64                         `json:"discount_amount"`
	CreatedAt       time.Time                       `json:"created_at"`
	UpdatedAt       time.Time                       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt                  `json:"deleted_at"`
}

func ToEventProductResponse(ep entities.EventProductEntity) EventProductResponse {
	return EventProductResponse{
		ID:              ep.ID,
		Event:           ep.Event.Name,
		Product:         productresponse.ToProductResponse(&ep.Product),
		DiscountPercent: ep.DiscountPercent,
		DiscountAmount:  ep.DiscountAmount,
		CreatedAt:       ep.CreatedAt,
		UpdatedAt:       ep.UpdatedAt,
		DeletedAt:       ep.DeletedAt,
	}
}
