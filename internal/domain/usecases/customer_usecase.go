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

	"github.com/google/uuid"
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

	existingPhone, err := u.repo.FindByPhoneCustomer(ctx, phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to check phone")
	}

	return existingPhone, nil
}

// RegisterCustomer implements services.CustomerService.
func (u *CustomerService) RegisterCustomer(ctx context.Context, req dto.RegisterCustomerRequest) (*entities.UserEntity, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("name is required")
	}
	if strings.TrimSpace(req.Username) == "" {
		return nil, errors.New("username is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return nil, errors.New("email is required")
	}
	if strings.TrimSpace(req.Phone) == "" {
		return nil, errors.New("phone is required")
	}
	if strings.TrimSpace(req.Gender) == "" {
		return nil, errors.New("gender is required")
	}
	role, err := u.repo.FindRoleByName(ctx, "customer")
	if err != nil {
		return nil, errors.New("failed to find role for customer")
	}
	existingUsername, err := u.repo.FindByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingUsername != nil {
		return nil, errors.New("username already exists")
	}

	existingEmail, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	existingPhone, err := u.repo.FindByPhone(ctx, req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingPhone != nil {
		return nil, errors.New("phone already exists")
	}

	newCustomer := &entities.UserEntity{
		ID:       uuid.NewString(),
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Gender:   &req.Gender,
		Phone:    &req.Phone,
		RoleID:   role.ID,
		Status:   0,
	}

	createdCustomer, err := u.repo.CreateCustomer(ctx, newCustomer)
	if err != nil {
		return nil, errors.New("failed to create customer")
	}

	return createdCustomer, nil
}
