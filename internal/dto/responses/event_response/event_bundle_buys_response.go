package eventresponse

import (
	"nusantara_service/internal/domain/entities"
	productresponse "nusantara_service/internal/dto/responses/product_response"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventBundleBuyResponse struct {
	ID        uuid.UUID                       `json:"id"`
	Event     string                          `json:"event"`
	Product   productresponse.ProductResponse `json:"product"`
	Quantity  int                             `json:"quantity"`
	CreatedAt time.Time                       `json:"created_at"`
	UpdatedAt time.Time                       `json:"updated_at"`
	DeletedAt gorm.DeletedAt                  `json:"deleted_at"`
}

func ToEvenBundleBuyResponse(ebb *entities.EventBundleBuyEntity) EventBundleBuyResponse {
	return EventBundleBuyResponse{
		ID:        ebb.ID,
		Event:     ebb.Event.Name,
		Product:   productresponse.ToProductResponse(&ebb.Product),
		Quantity:  ebb.Quantity,
		CreatedAt: ebb.CreatedAt,
		UpdatedAt: ebb.UpdatedAt,
		DeletedAt: ebb.DeletedAt,
	}
}
