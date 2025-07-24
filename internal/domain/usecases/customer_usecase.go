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
	digitsOnly := strings.TrimPrefix(normalized, "+")

	if len(digitsOnly) < 11 || len(digitsOnly) > 13 {
		return nil, response.NewCustomError(response.ErrBadRequest, "phone number must be 11 to 13 digits", 400)
	}

	if !utils.IsPhoneDigitsOnly(normalized) {
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

	err = configs.SetRedis(ctx, redisKey, otpCode, time.Minute*1)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to save OTP", 500)
	}

	err = twilio.SendWhatsAppOTP(normalizedPhone, otpCode)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to send OTP", 500)
	}

	return createdCustomer, nil
}

func (u *CustomerService) ResendCodeOTPVerify(ctx context.Context, req dto.ResendOTPRequest) error {
	if strings.TrimSpace(req.Phone) == "" {
		return response.NewCustomError(response.ErrBadRequest, "phone is required", 400)
	}

	normalizedPhone := utils.NormalizePhone(req.Phone)
	user, err := u.repo.FindByPhoneCustomer(ctx, normalizedPhone)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.NewCustomError(response.ErrNotFound, "phone number not registered", 404)
	}
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to check phone number", 500)
	}

	redisKey := fmt.Sprintf("otp:%s", *user.Phone)
	otpCode := otp.GenerateOTP(6)
	if err := configs.SetRedis(ctx, redisKey, otpCode, time.Minute*1); err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to store OTP", 500)
	}

	if err := twilio.SendWhatsAppOTP(*user.Phone, otpCode); err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to send OTP", 500)
	}

	return nil
}

// VerifyCodeOTP implements services.CustomerService.
func (u *CustomerService) VerifyCodeOTP(ctx context.Context, req dto.VerifyOTPRequest) error {
	if strings.TrimSpace(req.Phone) == "" || strings.TrimSpace(req.Code) == "" {
		return response.NewCustomError(response.ErrBadRequest, "phone and code are required", 400)
	}

	normalizedPhone := utils.NormalizePhone(req.Phone)

	user, err := u.repo.FindByPhoneCustomer(ctx, normalizedPhone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return response.NewCustomError(response.ErrInternal, "failed to check phone number", 500)
	}

	redisKey := fmt.Sprintf("otp:%s", normalizedPhone)
	storedCode, err := configs.GetRedis(ctx, redisKey)
	if err != nil {
		return response.NewCustomError(response.ErrUnauthorized, "OTP expired or invalid", 401)
	}

	if storedCode != req.Code {
		return response.NewCustomError(response.ErrUnauthorized, "invalid OTP code", 401)
	}

	_ = configs.DeleteRedis(ctx, redisKey)

	err = u.repo.UpdateStatusCustomer(ctx, user.ID, 1)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to update user status", 500)
	}

	return nil
}
