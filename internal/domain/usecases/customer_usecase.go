package usecases

import (
	"context"
	"errors"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/utils"
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
	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		return nil, errors.New("phone number is required")
	}

	normalized := utils.NormalizePhone(phone)

	if len(normalized) < 12 || len(normalized) > 13 {
		return nil, errors.New("phone number must be 12 or 13 digits")
	}

	if !utils.IsDigitsOnly(normalized) {
		return nil, errors.New("phone number must contain only digits")
	}

	existingPhone, err := u.repo.FindByPhone(ctx, phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to check phone")
	}

	return existingPhone, nil
}
