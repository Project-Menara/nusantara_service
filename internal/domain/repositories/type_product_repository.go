package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"

	"github.com/google/uuid"
)

type TypeProductRepository interface {
	Create(ctx context.Context, typeProduct *entities.TypeProductEntity) (*entities.TypeProductEntity, error)
	FindAll(ctx context.Context, limit, offset int) ([]*entities.TypeProductEntity, error)
	FindById(ctx context.Context, id uuid.UUID) (*entities.TypeProductEntity, error)
	Update(ctx context.Context, id uuid.UUID, data *entities.TypeProductEntity) (*entities.TypeProductEntity, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status int) error
	GetAllTypeProductCustomer(ctx context.Context) ([]*entities.TypeProductEntity, error)
	GetByIdTypeProductCustomer(ctx context.Context, id uuid.UUID) (*entities.TypeProductEntity, error)

	FindByName(ctx context.Context, name string) (*entities.TypeProductEntity, error)
	CountAll(ctx context.Context) (int, error)
	FindByUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error)

	FindAllWithSearch(ctx context.Context, limit, offset int, search string) ([]*entities.TypeProductEntity, error)
	CountAllWithSearch(ctx context.Context, search string) (int, error)
}
