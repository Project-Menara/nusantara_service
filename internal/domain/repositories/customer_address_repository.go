package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"

	"github.com/google/uuid"
)

type CustomerAddressRepository interface {
	Create(ctx context.Context, data *entities.CustomerAddressEntity) error
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.CustomerAddressEntity, error)
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.CustomerAddressEntity, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, data *entities.CustomerAddressEntity) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	FindCustomerID(ctx context.Context, userID string) (*entities.UserEntity, error)
	ClearDefaultByUser(ctx context.Context, userID uuid.UUID) error
	FindDefaultByUser(ctx context.Context, userID uuid.UUID) (*entities.CustomerAddressEntity, error)

	FindNearby(ctx context.Context, lat float64, lng float64, radiusKm float64) ([]*entities.ShopEntity, map[uuid.UUID]float64, error)

	SetDefaultByUser(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error
	SetFirstAsDefault(ctx context.Context, userID uuid.UUID) error
}
