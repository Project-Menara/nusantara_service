package usecases

import (
	"context"
	"errors"
	"fmt"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type authService struct {
	repo repositories.UserRepository
	rdb  *redis.Client
}

func NewAuthUsecase(repo repositories.UserRepository, rdb *redis.Client) services.AuthService {
	return &authService{repo: repo, rdb: rdb}
}

func (a *authService) Register(ctx context.Context, req dto.RegisterRequest) error {
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &entities.UserEntity{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashed,
		CreatedAt: time.Now(),
	}

	return a.repo.Create(ctx, user)
}

func (a *authService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	user, err := a.repo.FindByEmail(ctx, req.Email)
	if err != nil || !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", errors.New("invalid credentials")
	}
	return utils.GenerateToken(user.ID)
}

func (a *authService) Logout(ctx context.Context, token string) error {
	exp, err := utils.GetExpiryFromToken(token)
	if err != nil {
		return err
	}

	dur := time.Until(exp)
	key := fmt.Sprintf("blacklist:%s", token)
	return a.rdb.Set(ctx, key, "blacklisted", dur).Err()

}
