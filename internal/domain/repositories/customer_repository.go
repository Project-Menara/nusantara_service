package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
)

type CustomerRepository interface {
	FindByPhoneCustomer(ctx context.Context, phone string) (*entities.UserEntity, error)
	FindByUsername(ctx context.Context, username string) (*entities.UserEntity, error)
	FindByEmail(ctx context.Context, email string) (*entities.UserEntity, error)
	FindByPhone(ctx context.Context, phone string) (*entities.UserEntity, error)
	FindRoleByName(ctx context.Context, role string) (*entities.RoleEntity, error)
	CreateCustomer(ctx context.Context, user *entities.UserEntity) (*entities.UserEntity, error)
	UpdateStatusCustomer(ctx context.Context, userID string, status int) error
}
