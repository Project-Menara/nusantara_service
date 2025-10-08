package eventresponse

import (
	"nusantara_service/internal/domain/entities"
	productresponse "nusantara_service/internal/dto/responses/product_response"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventBundleRewardResponse struct {
	ID        uuid.UUID                       `json:"id"`
	Event     string                          `json:"event"`
	Product   productresponse.ProductResponse `json:"product"`
	Quantity  int                             `json:"quantity"`
	CreatedAt time.Time                       `json:"created_at"`
	UpdatedAt time.Time                       `json:"updated_at"`
	DeletedAt gorm.DeletedAt                  `json:"deleted_at"`
}

func ToEventBundleRewardResponse(ebr *entities.EventBundleRewardEntity) EventBundleRewardResponse {
	return EventBundleRewardResponse{
		ID:        ebr.ID,
		Event:     ebr.Event.Name,
		Product:   productresponse.ToProductResponse(&ebr.Product),
		Quantity:  ebr.Quantity,
		CreatedAt: ebr.CreatedAt,
		UpdatedAt: ebr.UpdatedAt,
		DeletedAt: ebr.DeletedAt,
	}
}
