package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
)

type UserRepository interface {
	CreateAdmin(ctx context.Context, user *entities.UserEntity) (*entities.UserEntity, error)
	FindExistUsername(ctx context.Context, username string) (*entities.UserEntity, error)
	FindByEmail(ctx context.Context, email string) (*entities.UserEntity, error)
	FindUserById(ctx context.Context, userId string) (*entities.UserEntity, error)
	ChangePassword(ctx context.Context, userId string, Updated *entities.UserEntity) (*entities.UserEntity, error)
}
