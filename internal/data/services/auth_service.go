package services

import (
	"context"
	"nusantara_service/internal/dto"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (string, error)
	Logout(ctx context.Context, token string) error
}
