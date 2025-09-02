package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"

	"github.com/google/uuid"
)

type CashierRepository interface {
	Create(ctx context.Context, data *entities.UserEntity) error
	FindByAll(ctx context.Context, offset, limit int, search string) ([]*entities.UserEntity, int, error)
	FindById(ctx context.Context, id uuid.UUID) (*entities.UserEntity, error)
	Update(ctx context.Context, id uuid.UUID, data *entities.UserEntity) error
	Delete(ctx context.Context, id uuid.UUID) error

	FindByEmail(ctx context.Context, email string) (*entities.UserEntity, error)
	FindByUsername(ctx context.Context, username string) (*entities.UserEntity, error)
}
