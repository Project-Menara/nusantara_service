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

func (u *CustomerService) CheckPhone(ctx context.Context, req dto.CheckPhoneRequest) (*response.CheckPhoneResult, error) {
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

	user, err := u.repo.FindByPhoneCustomer(ctx, normalized)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, response.NewCustomError(response.ErrInternal, "failed to check phone number", 500)
	}

	if user == nil {
		// Not registered
		return &response.CheckPhoneResult{
			Action: "register",
			User:   nil,
		}, nil
	}

	if user.Status == 0 {
		otpCode := otp.GenerateOTP(6)
		redisKey := fmt.Sprintf("otp:%s", normalized)

		err = configs.SetRedis(ctx, redisKey, otpCode, time.Minute*1)
		if err != nil {
			return nil, response.NewCustomError(response.ErrInternal, "failed to save OTP", 500)
		}

		err = twilio.SendWhatsAppOTP(normalized, otpCode)
		if err != nil {
			return nil, response.NewCustomError(response.ErrInternal, "failed to send OTP", 500)
		}
		// Belum verifikasi OTP
		return &response.CheckPhoneResult{
			Action: "verify_otp",
			User:   user,
		}, nil
	}

	if user.Status == 1 && user.Password == "" {
		otpCode := otp.GenerateOTP(6)
		redisKey := fmt.Sprintf("otp:%s", normalized)

		err = configs.SetRedis(ctx, redisKey, otpCode, time.Minute*1)
		if err != nil {
			return nil, response.NewCustomError(response.ErrInternal, "failed to save OTP", 500)
		}

		err = twilio.SendWhatsAppOTP(normalized, otpCode)
		if err != nil {
			return nil, response.NewCustomError(response.ErrInternal, "failed to send OTP", 500)
		}
		return &response.CheckPhoneResult{
			Action: "verify_otp_and_create_pin",
			User:   user,
		}, nil
	}

	// Siap login
	return &response.CheckPhoneResult{
		Action: "login",
		User:   user,
	}, nil
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

// NewPinCustomer implements services.CustomerService.
func (u *CustomerService) NewPinCustomer(ctx context.Context, req dto.NewPinRequest) error {
	if strings.TrimSpace(req.PIN) == "" || strings.TrimSpace(req.Phone) == "" {
		return response.NewCustomError(response.ErrBadRequest, "phone and pin required", 400)
	}

	normalizedPhone := utils.NormalizePhone(req.Phone)
	redisKey := fmt.Sprintf("newpin:%s", normalizedPhone)

	err := configs.SetRedis(ctx, redisKey, req.PIN, 15*time.Minute)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to save pin temporarily", 500)
	}

	return nil
}

// ConfirmPinCustomer implements services.CustomerService.
func (u *CustomerService) ConfirmPinCustomer(ctx context.Context, req dto.ConfirmPinRequest) (*entities.UserEntity, string, error) {
	normalizedPhone := utils.NormalizePhone(req.Phone)
	redisKey := fmt.Sprintf("newpin:%s", normalizedPhone)

	storedPIN, err := configs.GetRedis(ctx, redisKey)
	if err != nil {
		return nil, "", response.NewCustomError(response.ErrUnauthorized, "PIN expired or not set", 401)
	}

	if storedPIN != req.ConfirmPIN {
		return nil, "", response.NewCustomError(response.ErrUnauthorized, "PIN does not match", 401)
	}

	user, err := u.repo.FindByPhoneCustomer(ctx, normalizedPhone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, "", response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	hashedPin, err := utils.HashPassword(req.ConfirmPIN)
	if err != nil {
		return nil, "", response.NewCustomError(response.ErrInternal, "failed to hash pin", 500)
	}

	err = u.repo.UpdatePinCustomer(ctx, user.ID, hashedPin)
	if err != nil {
		return nil, "", response.NewCustomError(response.ErrInternal, "failed to save pin", 500)
	}

	_ = configs.DeleteRedis(ctx, redisKey)

	token, err := utils.GenerateToken(user.ID, user.Role.Name)
	if err != nil {
		return nil, "", response.NewCustomError(response.ErrInternal, "failed to generate token", 500)
	}

	redisTokenKey := fmt.Sprintf("customer_token:%s", user.ID)
	err = configs.SetRedis(ctx, redisTokenKey, token, 30*time.Minute)
	if err != nil {
		return nil, "", response.NewCustomError(response.ErrInternal, "failed to store session", 500)
	}

	return user, token, nil
}

// GetProfileCustomer implements services.CustomerService.
func (u *CustomerService) GetProfileCustomer(ctx context.Context, userId string, token string) (*entities.UserEntity, error) {
	redis_key := fmt.Sprintf("customer_token:%s", userId)
	storedToken, err := u.rdb.Get(ctx, redis_key).Result()
	if err != nil || storedToken != token {
		return nil, response.NewCustomError(response.ErrNotFound, "token not found", 404)
	}

	customer, err := u.repo.FindByUseIDCustomer(ctx, userId)
	if err != nil {
		return nil, response.NewCustomError(response.ErrNotFound, "user not found", 404)
	}

	return customer, nil
}

// LogoutCustomer implements services.CustomerService.
func (u *CustomerService) LogoutCustomer(ctx context.Context, userId, token string) error {
	expiry, err := utils.GetExpiryFromToken(token)
	if err != nil {
		return response.NewCustomError(response.ErrNotFound, "token expired not found", 404)
	}

	blackListKey := fmt.Sprintf("blacklist_customer:%s", token)
	err = u.rdb.Set(ctx, blackListKey, "blacklisted", time.Until(expiry)).Err()
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to blacklist token", 500)
	}

	redisKey := fmt.Sprintf("customer_token:%s", userId)
	err = u.rdb.Del(ctx, redisKey).Err()
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to delete session token", 500)
	}

	return nil

}

func (u *CustomerService) CheckTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	val, err := u.rdb.Get(ctx, fmt.Sprintf("blacklist_customer:%s", token)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil // Key does not exist, so not blacklisted
		}
		return false, err // Other actual Redis error
	}
	return val == "blacklisted", nil // Key exists, check its value

}
