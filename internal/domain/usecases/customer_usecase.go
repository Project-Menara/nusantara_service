package usecases

import (
	"context"
	"errors"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"nusantara_service/internal/dto"
	"strings"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CustomerService struct {
	repo repositories.CustomerRepository
	rdb  *redis.Client
}

func NewCustomerUsecase(repo repositories.CustomerRepository, rdb *redis.Client) services.CustomerService {
	return &CustomerService{repo: repo, rdb: rdb}
}

func (u *CustomerService) CheckPhone(ctx context.Context, req dto.CheckPhoneRequest) (*entities.UserEntity, error) {
	existingPhone, err := u.repo.FindByPhone(ctx, req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("phone not found")
	}
	if strings.TrimSpace(req.Phone) == "" {
		return nil, errors.New("phone number is required")
	}
	if len(req.Phone) < 12 || len(req.Phone) < 13 || len(req.Phone) > 13 {
		return nil, errors.New("phone number doesn't match")
	}

	return existingPhone, nil
}
