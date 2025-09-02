package services

import (
	"context"
	"nusantara_service/internal/domain/entities"
	dto "nusantara_service/internal/dto/request"
)

type UserService interface {

	//Super Admin
	RegisterAdmin(ctx context.Context, req dto.RegisterAdminRequest) (*entities.UserEntity, error)
	LoginAdmin(ctx context.Context, req dto.LoginAdminRequest) (string, error)
	GetProfile(ctx context.Context, userId string, token string) (*entities.UserEntity, error)
	LogoutAdmin(ctx context.Context, userId, token string) error
	CheckTokenBlacklisted(ctx context.Context, token string) (bool, error)
	ChangePasswordSuperAdmin(ctx context.Context, userId string, token string, req dto.ChangePasswordRequest) (*entities.UserEntity, error)
}
