package usecases

import (
	"context"
	"errors"
	"fmt"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/twilio"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	otp "nusantara_service/internal/utils/otp"
	"strings"
	"time"

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
		return nil, response.NewCustomError(response.ErrBadRequest, "phone is required", 400)
	}

	normalized := utils.NormalizePhone(phone)

	if len(normalized) < 12 || len(normalized) > 13 {
		return nil, response.NewCustomError(response.ErrBadRequest, "phone number must be 12 or 13 digits", 400)
	}

	if !utils.IsDigitsOnly(normalized) {
		return nil, response.NewCustomError(response.ErrBadRequest, "phone number must contain only digits", 400)
	}

	existingPhone, err := u.repo.FindByPhoneCustomer(ctx, normalized)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, response.NewCustomError(response.ErrInternal, "failed to check phone number", 500)
	}

	return existingPhone, nil
}

func (u *CustomerService) RegisterCustomer(ctx context.Context, req dto.RegisterCustomerRequest) (*entities.UserEntity, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "name is required", 400)
	}
	if strings.TrimSpace(req.Username) == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "username is required", 400)
	}
	if strings.TrimSpace(req.Email) == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "email is required", 400)
	}
	if strings.TrimSpace(req.Phone) == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "phone is required", 400)
	}
	if strings.TrimSpace(req.Gender) == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "gender is required", 400)
	}

	normalizedPhone := utils.NormalizePhone(req.Phone)

	role, err := u.repo.FindRoleByName(ctx, "customer")
	if err != nil {
		return nil, response.NewCustomError(response.ErrNotFound, "failed to find role for customer", 404)
	}

	if user, _ := u.repo.FindByUsername(ctx, req.Username); user != nil {
		return nil, response.NewCustomError(response.ErrExists, "username already exists", 409)
	}

	if user, _ := u.repo.FindByEmail(ctx, req.Email); user != nil {
		return nil, response.NewCustomError(response.ErrExists, "email already exists", 409)
	}

	if user, _ := u.repo.FindByPhone(ctx, normalizedPhone); user != nil {
		return nil, response.NewCustomError(response.ErrExists, "phone already exists", 409)
	}

	if !strings.HasSuffix(strings.ToLower(req.Email), "@gmail.com") {
		return nil, response.NewCustomError(response.ErrBadRequest, "only Gmail addresses are allowed", 400)
	}

	newCustomer := &entities.UserEntity{
		ID:       uuid.NewString(),
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Gender:   &req.Gender,
		Phone:    &normalizedPhone,
		RoleID:   role.ID,
		Status:   0,
	}

	createdCustomer, err := u.repo.CreateCustomer(ctx, newCustomer)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to create customer", 500)
	}

	otpCode := otp.GenerateOTP(6)
	redisKey := fmt.Sprintf("otp:%s", normalizedPhone)

	err = configs.SetRedis(ctx, redisKey, otpCode, time.Minute)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to save OTP", 500)
	}

	err = twilio.SendWhatsAppOTP(normalizedPhone, otpCode)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to send OTP", 500)
	}

	return createdCustomer, nil
}
