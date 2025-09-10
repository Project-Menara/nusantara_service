package customerresponse

import (
	"nusantara_service/internal/domain/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerAddressResponse struct {
	ID          uuid.UUID      `json:"id"`
	User        string         `json:"user"`
	Label       string         `json:"label"`
	AddressText string         `json:"address_text"`
	Lat         float64        `json:"lat"`
	Lng         float64        `json:"lang"`
	IsDefault   bool           `json:"is_default"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func ToCustomerAddressResponse(customerAddress entities.CustomerAddressEntity) CustomerAddressResponse {
	return CustomerAddressResponse{
		ID:          customerAddress.ID,
		User:        customerAddress.User.Name,
		Label:       customerAddress.Label,
		AddressText: customerAddress.AddressText,
		Lat:         customerAddress.Lat,
		Lng:         customerAddress.Lng,
		IsDefault:   customerAddress.IsDefault,
		CreatedAt:   customerAddress.CreatedAt,
		UpdatedAt:   customerAddress.UpdatedAt,
		DeletedAt:   customerAddress.DeletedAt,
	}
}
