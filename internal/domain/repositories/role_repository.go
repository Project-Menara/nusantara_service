package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"

	"github.com/google/uuid"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entities.RoleEntity) (*entities.RoleEntity, error)
	FindAll(ctx context.Context) ([]*entities.RoleEntity, error)
	FindById(ctx context.Context, Id uuid.UUID) (*entities.RoleEntity, error)
	Update(ctx context.Context, id uuid.UUID, updated *entities.RoleEntity) (*entities.RoleEntity, error)
	Delete(ctx context.Context, id uuid.UUID) error

	FindByName(ctx context.Context, name string) (*entities.RoleEntity, error)
}
