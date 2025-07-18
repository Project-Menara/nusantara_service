package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.UserEntity) error
	FindByEmail(ctx context.Context, email string) (*entities.UserEntity, error)
}
