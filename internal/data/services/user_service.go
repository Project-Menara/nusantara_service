package services

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/dto"
)

type UserService interface {
	RegisterAdmin(ctx context.Context, req dto.RegisterAdminRequest) (*entities.UserEntity, error)
	LoginAdmin(ctx context.Context, req dto.LoginAdminRequest) (string, error)
	GetProfile(ctx context.Context, userId string, token string) (*entities.UserEntity, error)
	LogoutAdmin(ctx context.Context, token string) error
	CheckTokenBlacklisted(ctx context.Context, token string) (bool, error)
}
