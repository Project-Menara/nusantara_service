package services

import (
	"context"
	"nusantara_service/internal/domain/entities"
	dto "nusantara_service/internal/dto/request"

	"github.com/google/uuid"
)

type CustomerAddressService interface {
	CreateAddres(ctx context.Context, custId uuid.UUID, req dto.CreateAddressRequest) error
	GetAllAddress(ctx context.Context, custID uuid.UUID) ([]*entities.CustomerAddressEntity, error)
	GetByIdAddress(ctx context.Context, id uuid.UUID, custID uuid.UUID) (*entities.CustomerAddressEntity, error)
	UpdateAddress(ctx context.Context, id uuid.UUID, custID uuid.UUID, req dto.UpdateAddressRequest) error
	Delete(ctx context.Context, id uuid.UUID, custID uuid.UUID) error

	SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID) error
	GetDefaultAddress(ctx context.Context, userID uuid.UUID) (*entities.CustomerAddressEntity, error)

	GetNearbyShops(ctx context.Context, lat, lng float64, maxDistance float64) ([]*entities.ShopEntity, map[uuid.UUID]float64, error)
}
