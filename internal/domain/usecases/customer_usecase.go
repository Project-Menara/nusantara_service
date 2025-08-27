package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/dataSources/rabbitmq"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	otp "nusantara_service/internal/utils/otp"
	"nusantara_service/internal/workers/consumer"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CustomerService struct {
	repo          repositories.CustomerRepository
	voucherRepo   repositories.VoucherRepository
	rdb           *redis.Client
	cloudinarySvc cloudinary.CloudinaryService
	db            *gorm.DB
}

func NewCustomerUsecase(repo repositories.CustomerRepository, rdb *redis.Client, cloudinarySvc *cloudinary.CloudinaryService, db *gorm.DB, voucherRepo repositories.VoucherRepository) services.CustomerService {
	return &CustomerService{repo: repo, rdb: rdb, cloudinarySvc: *cloudinarySvc, db: db, voucherRepo: voucherRepo}
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

	if !utils.IsPhoneDigitsOnly(phone) {
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

		err := configs.SetRedis(ctx, redisKey, otpCode, time.Minute*1)
		if err != nil {
			return nil, response.NewCustomError(response.ErrInternal, "failed to save OTP", 500)
		}

		// err = twilio.SendWhatsAppOTP(normalized, otpCode)
		// if err != nil {
		// 	return nil, response.NewCustomError(response.ErrInternal, err.Error(), 500)
		// }
		_ = rabbitmq.PublishToQueue("", "otp_queue", consumer.OTPPayload{
			Phone: normalized,
			Code:  otpCode,
		})

		ttl, err := u.rdb.TTL(ctx, redisKey).Result()
		if err != nil {
			return nil, err
		}
		// Belum verifikasi OTP
		return &response.CheckPhoneResult{
			Action: "verify_otp",
			User:   user,
			Ttl:    &ttl,
		}, nil
	}

	if user.Status == 1 && user.Password == "" {
		otpCode := otp.GenerateOTP(6)
		redisKey := fmt.Sprintf("otp:%s", normalized)

		err = configs.SetRedis(ctx, redisKey, otpCode, time.Minute*1)
		if err != nil {
			return nil, response.NewCustomError(response.ErrInternal, "failed to save OTP", 500)
		}

		// err = twilio.SendWhatsAppOTP(normalized, otpCode)
		// if err != nil {
		// 	return nil, response.NewCustomError(response.ErrInternal, "failed to send OTP", 500)
		// }
		_ = rabbitmq.PublishToQueue("", "otp_queue", consumer.OTPPayload{
			Phone: normalized,
			Code:  otpCode,
		})

		ttl, err := u.rdb.TTL(ctx, redisKey).Result()
		if err != nil {
			return nil, err
		}
		return &response.CheckPhoneResult{
			Action: "verify_otp_and_create_pin",
			User:   user,
			Ttl:    &ttl,
		}, nil
	}

	// Siap login
	return &response.CheckPhoneResult{
		Action: "login",
		User:   user,
	}, nil
}

func (u *CustomerService) RegisterCustomer(ctx context.Context, req dto.RegisterCustomerRequest) (*entities.UserEntity, time.Duration, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, 0, response.NewCustomError(response.ErrBadRequest, "name is required", 400)
	}
	if strings.TrimSpace(req.Username) == "" {
		return nil, 0, response.NewCustomError(response.ErrBadRequest, "username is required", 400)
	}
	if strings.TrimSpace(req.Email) == "" {
		return nil, 0, response.NewCustomError(response.ErrBadRequest, "email is required", 400)
	}
	if strings.TrimSpace(req.Phone) == "" {
		return nil, 0, response.NewCustomError(response.ErrBadRequest, "phone is required", 400)
	}
	if strings.TrimSpace(req.Gender) == "" {
		return nil, 0, response.NewCustomError(response.ErrBadRequest, "gender is required", 400)
	}

	normalizedPhone := utils.NormalizePhone(req.Phone)

	role, err := u.repo.FindRoleByName(ctx, "customer")
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrNotFound, "failed to find role for customer", 404)
	}

	if user, _ := u.repo.FindByUsername(ctx, req.Username); user != nil {
		return nil, 0, response.NewCustomError(response.ErrExists, "username already exists", 409)
	}

	if user, _ := u.repo.FindByEmail(ctx, req.Email); user != nil {
		return nil, 0, response.NewCustomError(response.ErrExists, "email already exists", 409)
	}

	if user, _ := u.repo.FindByPhone(ctx, normalizedPhone); user != nil {
		return nil, 0, response.NewCustomError(response.ErrExists, "phone already exists", 409)
	}

	if !strings.HasSuffix(strings.ToLower(req.Email), "@gmail.com") {
		return nil, 0, response.NewCustomError(response.ErrBadRequest, "only Gmail addresses are allowed", 400)
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
		return nil, 0, response.NewCustomError(response.ErrInternal, "failed to create customer", 500)
	}

	customerPoint := &entities.UserPointEntity{
		ID:          uuid.New(),
		UserID:      uuid.MustParse(createdCustomer.ID),
		TotalPoints: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_ = u.repo.CreateCustomerPoint(ctx, customerPoint)

	otpCode := otp.GenerateOTP(6)
	redisKey := fmt.Sprintf("otp:%s", normalizedPhone)

	err = configs.SetRedis(ctx, redisKey, otpCode, time.Minute*1)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrInternal, "failed to save OTP", 500)
	}

	// err = twilio.SendWhatsAppOTP(normalizedPhone, otpCode)
	// if err != nil {
	// 	return nil, 0, response.NewCustomError(response.ErrInternal, "failed to send OTP", 500)
	// }
	_ = rabbitmq.PublishToQueue("", "otp_queue", consumer.OTPPayload{
		Phone: normalizedPhone,
		Code:  otpCode,
	})

	ttl, err := u.rdb.TTL(ctx, redisKey).Result()
	if err != nil {
		return nil, 0, err
	}

	return createdCustomer, time.Duration(ttl.Seconds()), nil
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

	err = rabbitmq.PublishToQueue("", "otp_queue", consumer.OTPPayload{
		Phone: *user.Phone,
		Code:  otpCode,
	})
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to initiate OTP resend", 500)
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

	_ = rabbitmq.PublishToQueue("verify_otp", "verified_queue", consumer.VerifiedPayload{
		Phone: normalizedPhone,
	})

	return nil
}

// NewPinCustomer implements services.CustomerService.
func (u *CustomerService) NewPinCustomer(ctx context.Context, req dto.NewPinRequest) error {
	pin := strings.TrimSpace(req.PIN)
	phone := strings.TrimSpace(req.Phone)

	if pin == "" || phone == "" {
		return response.NewCustomError(response.ErrBadRequest, "phone and pin required", 400)
	}
	if !utils.IsDigitsOnly(pin) {
		return response.NewCustomError(response.ErrBadRequest, "PIN must contain only digits", 400)
	}
	if len(pin) != 6 {
		return response.NewCustomError(response.ErrBadRequest, "PIN must be 6 digits", 400)
	}

	normalizedPhone := utils.NormalizePhone(phone)
	redisKey := fmt.Sprintf("newpin:%s", normalizedPhone)

	err := configs.SetRedis(ctx, redisKey, pin, 15*time.Minute)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to save pin temporarily", 500)
	}

	return nil
}

// ConfirmPinCustomer implements services.CustomerService.
func (u *CustomerService) ConfirmPinCustomer(ctx context.Context, req dto.ConfirmPinRequest) (*entities.UserEntity, string, error) {
	pin := strings.TrimSpace(req.ConfirmPIN)
	phone := strings.TrimSpace(req.Phone)

	if pin == "" || phone == "" {
		return nil, "", response.NewCustomError(response.ErrBadRequest, "phone and pin required", 400)
	}
	if !utils.IsDigitsOnly(pin) {
		return nil, "", response.NewCustomError(response.ErrBadRequest, "Confirmation PIN must contain only digits", 400)
	}
	if len(pin) != 6 {
		return nil, "", response.NewCustomError(response.ErrBadRequest, "Confirmation PIN must be 6 digits", 400)
	}

	normalizedPhone := utils.NormalizePhone(phone)
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

	hashedPin, err := utils.HashPassword(pin)
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

// LoginCustomer implements services.CustomerService.
func (u *CustomerService) LoginCustomer(ctx context.Context, req dto.LoginCustomerRequest) (string, error) {
	if strings.TrimSpace(req.Phone) == "" {
		return "", response.NewCustomError(response.ErrBadRequest, "phone is required", 400)
	}
	if strings.TrimSpace(req.Pin) == "" {
		return "", response.NewCustomError(response.ErrBadRequest, "pin is required", 400)
	}

	if !utils.IsDigitsOnly(req.Pin) {
		return "", response.NewCustomError(response.ErrBadRequest, "Confirmation PIN must contain only digits", 400)
	}
	if len(req.Pin) != 6 {
		return "", response.NewCustomError(response.ErrBadRequest, "Confirmation PIN must be 6 digits", 400)
	}

	if req.Pin == "" || req.Phone == "" {
		return "", response.NewCustomError(response.ErrBadRequest, "phone and pin required", 400)
	}
	if !utils.IsDigitsOnly(req.Pin) {
		return "", response.NewCustomError(response.ErrBadRequest, "PIN must contain only digits", 400)
	}
	if len(req.Pin) != 6 {
		return "", response.NewCustomError(response.ErrBadRequest, "PIN must be 6 digits", 400)
	}

	loginKey := fmt.Sprintf("login_attempts_customer:%s", req.Phone)

	const maxAttempts = 5
	const cooldownDuration = 1 * time.Minute

	attempts, err := u.rdb.Get(ctx, loginKey).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", response.NewCustomError(response.ErrInternal, "failed to get login attempts from redis", 500)
	}

	ttl, err := u.rdb.TTL(ctx, loginKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", response.NewCustomError(response.ErrInternal, "failed to get TTL from redis", 500)
	}

	if attempts >= maxAttempts {
		if ttl > 0 {
			return "", &response.CooldownError{
				Message:    "Too many login attempts. Please try again later.",
				RetryAfter: ttl,
			}
		}
	}

	normalizedPhone := utils.NormalizePhone(req.Phone)

	customer, err := u.repo.FindByPhoneCustomer(ctx, normalizedPhone)
	if err != nil {
		currentAttempts, err := u.rdb.Incr(ctx, loginKey).Result()
		if err != nil {
			return "", response.NewCustomError(response.ErrInternal, "failed to increment login attempts", 500)
		}
		u.rdb.Expire(ctx, loginKey, cooldownDuration)
		if int(currentAttempts) >= maxAttempts {
			ttl, _ := u.rdb.TTL(ctx, loginKey).Result()
			if ttl <= 0 {
				ttl = cooldownDuration
			}
			return "", &response.CooldownError{
				Message:    "Too many login attempts. Please try again later.",
				RetryAfter: ttl,
			}
		}

		remainingAttempts := maxAttempts - int(currentAttempts)
		if remainingAttempts < 0 {
			remainingAttempts = 0
		}

		return "", &response.LoginAttemptError{
			Message:           "phone incorrect",
			RemainingAttempts: remainingAttempts,
		}
	}

	isPin := utils.CheckPasswordHash(req.Pin, customer.Password)
	if !isPin {
		currentAttempts, err := u.rdb.Incr(ctx, loginKey).Result()
		if err != nil {
			return "", fmt.Errorf("failed to increment login attempts: %w", err)
		}
		u.rdb.Expire(ctx, loginKey, cooldownDuration)

		if int(currentAttempts) >= maxAttempts {
			ttl, _ := u.rdb.TTL(ctx, loginKey).Result()
			if ttl <= 0 {
				ttl = cooldownDuration
			}
			return "", &response.CooldownError{
				Message:    "Too many login attempts. Please try again later.",
				RetryAfter: ttl,
			}
		}

		remainingAttempts := maxAttempts - int(currentAttempts)
		if remainingAttempts < 0 {
			remainingAttempts = 0
		}

		return "", &response.LoginAttemptError{
			Message:           "pin incorrect",
			RemainingAttempts: remainingAttempts,
		}
	}

	u.rdb.Del(ctx, loginKey)

	token, err := utils.GenerateToken(customer.ID, customer.Role.Name)
	if err != nil {
		return "", response.NewCustomError(response.ErrInternal, "failed to generate token", 500)
	}

	expiry, err := utils.GetExpiryFromToken(token)
	if err != nil {
		return "", response.NewCustomError(response.ErrInternal, "failed to get token expiry", 500)
	}

	redisKey := fmt.Sprintf("customer_token:%s", customer.ID)
	err = u.rdb.Set(ctx, redisKey, token, time.Until(expiry)).Err()
	if err != nil {
		return "", response.NewCustomError(response.ErrInternal, "failed to store token in redis", 500)
	}

	return token, nil
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

// UpdateProfileCustomer implements services.CustomerService.
func (u *CustomerService) UpdateProfileCustomer(ctx context.Context, userId string, req dto.UpdateCustomerRequest, photoFileHeader *multipart.FileHeader) (*entities.UserEntity, error) {
	user, err := u.repo.FindByUseIDCustomer(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Username != "" {
		existingUser, _ := u.repo.FindByUsername(ctx, req.Username)
		if existingUser != nil && existingUser.ID != userId {
			return nil, response.NewCustomError(response.ErrExists, "username already exists", 409)
		}
		user.Username = req.Username
	}

	if req.Email != "" {
		existingUser, _ := u.repo.FindByEmail(ctx, req.Email)
		if existingUser != nil && existingUser.ID != userId {
			return nil, response.NewCustomError(response.ErrExists, "email already exists", 409)
		}
		if !strings.HasSuffix(strings.ToLower(req.Email), "@gmail.com") {
			return nil, response.NewCustomError(response.ErrBadRequest, "only Gmail addresses are allowed", 400)
		}
		user.Email = req.Email
	}

	if req.Gender != "" {
		user.Gender = &req.Gender
	}
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			return nil, response.NewCustomError(response.ErrBadRequest, "invalid date of birth format, use YYYY-MM-DD", 400)
		}
		user.DateOfBirth = &dob
	}

	if photoFileHeader != nil {
		if user.Photo != nil && *user.Photo != "" {
			publicID := utils.ExtractPublicIDFromCloudinaryURL(*user.Photo)
			if publicID != "" {
				if err := u.cloudinarySvc.DestroyImage(ctx, publicID); err != nil {
					fmt.Printf("Warning: Failed to delete old image from Cloudinary: %v\n", err)
				}
			}
		}

		folder := fmt.Sprintf("nusantara_service/customer_profiles/%s", user.Username)
		filename := fmt.Sprintf("profile_%s", userId)
		photoURL, err := u.cloudinarySvc.UploadImage(ctx, photoFileHeader, folder, filename)
		if err != nil {
			return nil, response.NewCustomError(response.ErrInternal, "failed to upload photo to Cloudinary", 500)
		}

		user.Photo = &photoURL.URL

	}

	updatedUser, err := u.repo.UpdateCustomer(ctx, userId, user)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to update customer profile", 500)
	}

	return updatedUser, nil

}

// VerifyPINCustomer implements services.CustomerService.
func (u *CustomerService) VerifyPINCustomer(ctx context.Context, userId string, req dto.VerifyPINCustomerRequest) error {
	if strings.TrimSpace(req.PIN) == "" {
		return response.NewCustomError(response.ErrBadRequest, "pin is required", 400)
	}

	if !utils.IsDigitsOnly(req.PIN) {
		return response.NewCustomError(response.ErrBadRequest, "Confirmation PIN must contain only digits", 400)
	}
	if len(req.PIN) != 6 {
		return response.NewCustomError(response.ErrBadRequest, "Confirmation PIN must be 6 digits", 400)
	}

	pinKey := fmt.Sprintf("pin_attempts:%s", userId)
	const maxAttempts = 5
	const cooldownDuration = 1 * time.Minute

	attempts, err := u.rdb.Get(ctx, pinKey).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return response.NewCustomError(response.ErrInternal, "failed to pinkey attempts from redis", 500)
	}

	ttl, err := u.rdb.TTL(ctx, pinKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return response.NewCustomError(response.ErrInternal, "failed to get TTL from redis", 500)
	}

	if attempts >= maxAttempts {
		if ttl > 0 {
			return &response.CooldownError{
				Message:    "Too many pin attempts. Please try again later.",
				RetryAfter: ttl,
			}
		}
	}

	customer, err := u.repo.FindByUseIDCustomer(ctx, userId)
	if err != nil {
		currentAttempts, err := u.rdb.Incr(ctx, pinKey).Result()
		if err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to increment pin attempts", 500)
		}

		u.rdb.Expire(ctx, pinKey, cooldownDuration)
		if int(currentAttempts) >= maxAttempts {
			ttl, _ := u.rdb.TTL(ctx, pinKey).Result()
			if ttl <= 0 {
				ttl = cooldownDuration
			}
			return &response.CooldownError{
				Message:    "Too many pin attempts. Please try again later.",
				RetryAfter: ttl,
			}
		}

		remainingAttempts := maxAttempts - int(currentAttempts)
		if remainingAttempts < 0 {
			remainingAttempts = 0
		}

		return &response.VerifyPINAttemptError{
			Message:           "user not permission",
			RemainingAttempts: remainingAttempts,
		}
	}

	isPin := utils.CheckPasswordHash(req.PIN, customer.Password)
	if !isPin {
		currentAttempts, err := u.rdb.Incr(ctx, pinKey).Result()
		if err != nil {
			return fmt.Errorf("failed to increment pin attempts: %w", err)
		}
		u.rdb.Expire(ctx, pinKey, cooldownDuration)

		if int(currentAttempts) >= maxAttempts {
			ttl, _ := u.rdb.TTL(ctx, pinKey).Result()
			if ttl <= 0 {
				ttl = cooldownDuration
			}
			return &response.CooldownError{
				Message:    "Too many pin attempts. Please try again later.",
				RetryAfter: ttl,
			}
		}

		remainingAttempts := maxAttempts - int(currentAttempts)
		if remainingAttempts < 0 {
			remainingAttempts = 0
		}

		return &response.VerifyPINAttemptError{
			Message:           "pin incorrect",
			RemainingAttempts: remainingAttempts,
		}
	}

	u.rdb.Del(ctx, pinKey)

	return nil
}

// NewPINCustomer implements services.CustomerService.
func (u *CustomerService) NewPINCustomer(ctx context.Context, userId string, req dto.NewPINCustomer) error {
	NewPin := strings.TrimSpace(req.NewPIN)

	if NewPin == "" {
		return response.NewCustomError(response.ErrBadRequest, "new pin required", 400)
	}

	if !utils.IsDigitsOnly(NewPin) {
		return response.NewCustomError(response.ErrBadRequest, "PIN must contain only digits", 400)
	}
	if len(NewPin) != 6 {
		return response.NewCustomError(response.ErrBadRequest, "PIN must be 6 digits", 400)
	}

	redisKey := fmt.Sprintf("new_pin_customer:%s", userId)
	err := configs.SetRedis(ctx, redisKey, NewPin, 15*time.Minute)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to save pin", 500)
	}

	return nil
}

// ConfirmationPINCustomerUpdate implements services.CustomerService.
func (u *CustomerService) ConfirmationPINCustomerUpdate(ctx context.Context, userId string, req dto.ConfirmNewPINCustomer) (*entities.UserEntity, error) {
	pin := strings.TrimSpace(req.ConfirmPIN)
	if pin == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "confirmation pin required", 400)
	}

	if !utils.IsDigitsOnly(pin) {
		return nil, response.NewCustomError(response.ErrBadRequest, "PIN must contain only digits", 400)
	}
	if len(pin) != 6 {
		return nil, response.NewCustomError(response.ErrBadRequest, "PIN must be 6 digits", 400)
	}

	redisKey := fmt.Sprintf("new_pin_customer:%s", userId)
	storedPin, err := configs.GetRedis(ctx, redisKey)
	if err != nil {
		return nil, response.NewCustomError(response.ErrUnauthorized, "PIN expired or not set", 401)
	}

	if storedPin != req.ConfirmPIN {
		return nil, response.NewCustomError(response.ErrUnauthorized, "PIN does not match", 401)
	}

	customer, err := u.repo.FindByUseIDCustomer(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	hashedPin, err := utils.HashPassword(pin)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to hash pin", 500)
	}

	err = u.repo.UpdatePinCustomer(ctx, userId, hashedPin)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to save pin", 500)
	}

	_ = configs.DeleteRedis(ctx, redisKey)

	return customer, nil
}

// NewPhoneCustomer implements services.CustomerService.
func (u *CustomerService) NewPhoneCustomer(ctx context.Context, userId string, req dto.NewPhoneCustomerRequest) error {
	NewPhone := strings.TrimSpace(req.Phone)

	if NewPhone == "" {
		return response.NewCustomError(response.ErrBadRequest, "New Phone required", 400)
	}

	normalized := utils.NormalizePhone(NewPhone)
	digitsOnly := strings.TrimPrefix(normalized, "+")

	if len(digitsOnly) < 11 || len(digitsOnly) > 13 {
		return response.NewCustomError(response.ErrBadRequest, "phone number must be 11 to 13 digits", 400)
	}

	if !utils.IsPhoneDigitsOnly(normalized) {
		return response.NewCustomError(response.ErrBadRequest, "phone number must contain only digits", 400)
	}

	_, err := u.repo.FindByUseIDCustomer(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	redisKey := fmt.Sprintf("otp_new_phone:%s", normalized)
	otpCode := otp.GenerateOTP(6)
	if err := configs.SetRedis(ctx, redisKey, otpCode, time.Minute*1); err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to store OTP", 500)
	}

	err = rabbitmq.PublishToQueue("", "otp_queue", consumer.OTPPayload{
		Phone: normalized,
		Code:  otpCode,
	})
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to initiate OTP resend", 500)
	}

	return nil
}

// VerifyCodeOTPCustomerUpdate implements services.CustomerService.
func (u *CustomerService) VerifyCodeOTPCustomerUpdate(ctx context.Context, userId string, req dto.VerifyOTPCustomerUpdateRequest) (*entities.UserEntity, error) {
	if strings.TrimSpace(req.Phone) == "" || strings.TrimSpace(req.Code) == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "phone and code are required", 400)
	}

	normalizedPhone := utils.NormalizePhone(req.Phone)

	customer, err := u.repo.FindByUseIDCustomer(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	redisKey := fmt.Sprintf("otp_new_phone:%s", normalizedPhone)
	storeCode, err := configs.GetRedis(ctx, redisKey)
	if err != nil {
		return nil, response.NewCustomError(response.ErrUnauthorized, "OTP expired or invalid", 401)
	}

	if storeCode != req.Code {
		return nil, response.NewCustomError(response.ErrUnauthorized, "invalid OTP code", 401)
	}

	_ = configs.DeleteRedis(ctx, redisKey)

	customer.Phone = &req.Phone

	data, err := u.repo.ChangePhoneCustomer(ctx, userId, customer)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "failed to update user phone", 500)
	}

	_ = rabbitmq.PublishToQueue("verify_otp", "verified_queue", consumer.VerifiedPayload{
		Phone: normalizedPhone,
	})

	return data, nil
}

func (u *CustomerService) FogotPIN(ctx context.Context, req dto.ForgotPINRequest) (string, error) {
	if strings.TrimSpace(req.Phone) == "" {
		return "", response.NewCustomError(response.ErrBadRequest, "phone is required", 400)
	}

	_, err := u.repo.FindByPhoneCustomer(ctx, req.Phone)
	if err != nil {
		return "", response.NewCustomError(response.ErrNotFound, "user not found", 500)
	}
	token := uuid.NewString()
	normalized := utils.NormalizePhone(req.Phone)

	deepLink := fmt.Sprintf("https://nusantara-oleh-oleh.com/reset-pin?token=%s", token)

	message := fmt.Sprintf("Klik link berikut untuk reset PIN: \n%s\nBerlaku selama 2 menit", deepLink)

	err = configs.SetRedis(ctx, "reset_pin:"+token, message, time.Minute*2)
	if err != nil {
		return "", response.NewCustomError(response.ErrInternal, "failed to store token", 404)
	}

	// err = twilio.SendWhatsAppMessage(normalized, message)
	// if err != nil {
	// 	return "", response.NewCustomError(response.ErrInternal, "failed to send message link", 500)
	// }
	_ = rabbitmq.PublishToQueue("forgot_pin", rabbitmq.LinkForgotPINQueueName, consumer.LinkForgotPINPayload{
		Phone: normalized,
		Link:  message,
	})

	return token, nil
}

// ValidateTokenForgotPIN implements services.CustomerService.
func (u *CustomerService) ValidateTokenForgotPIN(ctx context.Context, token string) (string, error) {
	if strings.TrimSpace(token) == "" {
		return "", response.NewCustomError(response.ErrBadRequest, "token is required", 400)
	}

	redisKey := "reset_pin:" + token
	phone, err := configs.GetRedis(ctx, redisKey)
	if err != nil {
		return "", response.NewCustomError(response.ErrUnauthorized, "invalid or expired token", 401)
	}

	return phone, nil
}

// CreateNewPIN implements services.CustomerService.
func (u *CustomerService) CreateNewPIN(ctx context.Context, token string, req dto.CreateNewPinRequest) error {
	newPin := strings.TrimSpace(req.PIN)
	phone := strings.TrimSpace(req.Phone)

	if newPin == "" {
		return response.NewCustomError(response.ErrBadRequest, "pin is required", 400)
	}
	if phone == "" {
		return response.NewCustomError(response.ErrBadRequest, "phone is required", 400)
	}

	if len(newPin) != 6 {
		return response.NewCustomError(response.ErrBadRequest, "PIN must be 6 digits", 400)
	}

	user, err := u.repo.FindByPhoneCustomer(ctx, req.Phone)
	if err != nil {
		return response.NewCustomError(response.ErrNotFound, "user not found", 500)
	}

	isOldPin := utils.CheckPasswordHash(req.PIN, user.Password)
	if isOldPin {
		return response.NewCustomError(response.ErrBadRequest, "You've input the old pin. Please input a new pin", 500)
	}

	normalize := utils.NormalizePhone(phone)
	redisKey := fmt.Sprintf("new_pin_forgot:%s", normalize)

	err = configs.SetRedis(ctx, redisKey, newPin, time.Minute*15)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to save pin", 500)
	}

	return nil
}

// CreateNewConfirmPIN implements services.CustomerService.
func (u *CustomerService) CreateNewConfirmPIN(ctx context.Context, tokenPIN string, req dto.CreateConfirmPinRequest) (*entities.UserEntity, string, error) {
	confirmPin := strings.TrimSpace(req.ConfirmPIN)
	phone := strings.TrimSpace(req.Phone)

	if confirmPin == "" {
		return nil, "", response.NewCustomError(response.ErrBadRequest, "phone is required", 400)
	}
	if phone == "" {
		return nil, "", response.NewCustomError(response.ErrBadRequest, "phone is required", 400)
	}
	normalize := utils.NormalizePhone(phone)
	redisKey := fmt.Sprintf("new_pin_forgot:%s", normalize)

	storedPIN, err := configs.GetRedis(ctx, redisKey)
	if err != nil {
		return nil, "", response.NewCustomError(response.ErrUnauthorized, "PIN expired or not set", 401)
	}

	if storedPIN != confirmPin {
		return nil, "", response.NewCustomError(response.ErrUnauthorized, "PIN doesn't match", 401)
	}

	user, err := u.repo.FindByPhoneCustomer(ctx, normalize)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, "", response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	hashedPin, err := utils.HashPassword(confirmPin)
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

// ClaimVoucherCustomer implements services.CustomerService.
func (u *CustomerService) ClaimVoucherCustomer(ctx context.Context, customerID uuid.UUID, voucherID uuid.UUID) (*entities.UserVoucherEntity, error) {
	tx := u.db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	voucher, err := u.voucherRepo.GetByIdVoucherCustomer(ctx, voucherID)
	if err != nil {
		tx.Rollback()
		return nil, response.NewCustomError(response.ErrNotFound, "voucher not found", 404)
	}

	now := time.Now()
	if now.After(voucher.EndDate) {
		tx.Rollback()
		return nil, response.NewCustomError(response.ErrBadRequest, "voucher expired", 400)
	}

	if voucher.ClaimedCount >= voucher.Quota {
		tx.Rollback()
		return nil, response.NewCustomError(response.ErrBadRequest, "voucher not available", 400)
	}

	claimed, err := u.voucherRepo.CheckUserVoucherClaimed(ctx, customerID, voucherID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if claimed {
		tx.Rollback()
		return nil, response.NewCustomError(response.ErrInternal, "voucher already claimed", 500)
	}

	if voucher.PointCost > 0 {
		customerPoint, err := u.repo.GetCustomerPoint(ctx, customerID)
		if err != nil {
			tx.Rollback()
			return nil, response.NewCustomError(response.ErrNotFound, "points data not found", 404)
		}
		if customerPoint.TotalPoints < voucher.PointCost {
			tx.Rollback()
			return nil, response.NewCustomError(response.ErrBadRequest, "not enough points to claim voucher", 400)
		}

		if err := u.repo.UpdateCustomerPoint(ctx, customerID, voucher.PointCost); err != nil {
			tx.Rollback()
			return nil, err
		}

		history := entities.UserPointHistoriesEntity{
			ID:          uuid.New(),
			UserID:      customerID,
			PointType:   "exchange",
			Source:      "voucher",
			SourceId:    voucher.ID.String(),
			Points:      voucher.PointCost,
			Direction:   "out",
			Description: "Claim voucher " + voucher.Code,
		}

		if err := u.repo.CreatePointHistory(ctx, &history); err != nil {
			tx.Rollback()
			return nil, response.NewCustomError(response.ErrInternal, "failed to create point history", 500)
		}
	}

	detailVoucher := entities.UserVoucherDetailEntity{
		ID:                uuid.New(),
		VoucherCode:       voucher.Code,
		DiscountType:      voucher.DiscountType,
		DiscountAmount:    voucher.DiscountAmount,
		DiscountPercent:   voucher.DiscountPercent,
		MinPurchaseAmount: voucher.MinimumSpend,
		ValidFrom:         voucher.StartDate,
		ValidUntil:        voucher.EndDate,
		Description:       voucher.Description,
	}

	if _, err := u.repo.AddDetailVoucher(ctx, &detailVoucher); err != nil {
		tx.Rollback()
		return nil, response.NewCustomError(response.ErrInternal, "failed to add voucher detail", 500)
	}

	userVoucher := entities.UserVoucherEntity{
		ID:        uuid.New(),
		UserID:    customerID,
		VoucherID: voucher.ID,
		DetailID:  detailVoucher.ID,
		IsUsed:    false,
	}

	created, err := u.repo.AddVoucher(ctx, &userVoucher)
	if err != nil {
		tx.Rollback()
		return nil, response.NewCustomError(response.ErrInternal, "failed to add user voucher", 500)
	}

	if err := u.voucherRepo.IncreaseVoucherClaimedCount(ctx, voucherID); err != nil {
		tx.Rollback()
		return nil, response.NewCustomError(response.ErrInternal, "failed to increase voucher claimed count", 500)
	}

	u.InvalidateCustomerCache(ctx)

	return created, tx.Commit().Error

}

func (v *CustomerService) InvalidateCustomerCache(ctx context.Context) {
	cp := v.rdb.Scan(ctx, 0, "customer_point:*", 0).Iterator()
	for cp.Next(ctx) {
		v.rdb.Del(ctx, cp.Val())
	}
	cph := v.rdb.Scan(ctx, 0, "customer_point_history:*", 0).Iterator()
	for cph.Next(ctx) {
		v.rdb.Del(ctx, cph.Val())
	}
	vc := v.rdb.Scan(ctx, 0, "customer_vouchers_claimed:*", 0).Iterator()
	for vc.Next(ctx) {
		v.rdb.Del(ctx, vc.Val())
	}
}
func (v *CustomerService) InvalidateVoucherCache(ctx context.Context) {
	vca := v.rdb.Scan(ctx, 0, "vouchers:*", 0).Iterator()
	for vca.Next(ctx) {
		v.rdb.Del(ctx, vca.Val())
	}
	vcai := v.rdb.Scan(ctx, 0, "voucher:*", 0).Iterator()
	for vcai.Next(ctx) {
		v.rdb.Del(ctx, vcai.Val())
	}
}

// GetCustomerPoint implements services.CustomerService.
func (u *CustomerService) GetCustomerPoint(ctx context.Context, customerID uuid.UUID) (*entities.UserPointEntity, int, *time.Time, error) {
	cacheKey := fmt.Sprintf("customer_point:%s", customerID.String())
	cached, err := configs.GetRedis(ctx, cacheKey)
	if err == nil && cached != "" {
		var userPoint entities.UserPointEntity
		if err := json.Unmarshal([]byte(cached), &userPoint); err == nil {
			return &userPoint, 0, nil, nil
		}
	}

	// Ambil data awal
	customerPoint, totalExpired, expiredDate, err := u.repo.FindUserPoint(ctx, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil, response.NewCustomError(response.ErrNotFound, "customer point not found", 404)
		}
		return nil, 0, nil, err
	}

	// Kalau ada yang expired
	if expiredDate != nil && time.Now().After(*expiredDate) && totalExpired > 0 {
		// Update jadi expired
		err = u.repo.MarkPointsAsExpired(ctx, customerID, *expiredDate)
		if err != nil {
			return nil, 0, nil, err
		}

		// Kurangi total points di tabel utama
		err = u.repo.DecreaseTotalPoints(ctx, customerID, totalExpired)
		if err != nil {
			return nil, 0, nil, err
		}

		// Ambil ulang data setelah update
		customerPoint, totalExpired, expiredDate, err = u.repo.FindUserPoint(ctx, customerID)
		if err != nil {
			return nil, 0, nil, err
		}
	}

	// Cache hasil
	dataCache, _ := json.Marshal(customerPoint)
	_ = configs.SetRedis(ctx, cacheKey, dataCache, time.Minute*30)

	u.InvalidateVoucherCache(ctx)
	u.InvalidateCustomerCache(ctx)

	return customerPoint, totalExpired, expiredDate, nil
}

// GetCustomerPointHistory implements services.CustomerService.
func (u *CustomerService) GetCustomerPointHistory(ctx context.Context, customerID uuid.UUID) ([]*entities.UserPointHistoriesEntity, error) {
	cacheKey := fmt.Sprintf("customer_point_history:%s", customerID.String())
	cached, err := configs.GetRedis(ctx, cacheKey)
	if err == nil && cached != "" {
		var pointHistories []*entities.UserPointHistoriesEntity
		if err := json.Unmarshal([]byte(cached), &pointHistories); err == nil {
			return pointHistories, nil
		}
	}

	pointHistories, err := u.repo.FindUserPointHistory(ctx, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "customer point history not found", 404)
		}
		return nil, err
	}

	dataCache, _ := json.Marshal(pointHistories)
	_ = configs.SetRedis(ctx, cacheKey, dataCache, time.Minute*30)

	u.InvalidateVoucherCache(ctx)
	u.InvalidateCustomerCache(ctx)
	return pointHistories, nil
}

// GetCustomerVouchersClaimed implements services.CustomerService.
func (u *CustomerService) GetCustomerVouchersClaimed(ctx context.Context, customerID uuid.UUID) ([]*entities.UserVoucherEntity, error) {
	cacheKey := fmt.Sprintf("customer_vouchers_claimed:%s", customerID.String())
	cached, err := configs.GetRedis(ctx, cacheKey)
	if err == nil && cached != "" {
		var userVouchers []*entities.UserVoucherEntity
		if err := json.Unmarshal([]byte(cached), &userVouchers); err == nil {
			return userVouchers, nil
		}
	}

	userVouchers, err := u.repo.FindUserVoucherClaimed(ctx, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "customer vouchers claimed not found", 404)
		}
		return nil, err
	}

	dataCache, _ := json.Marshal(userVouchers)
	_ = configs.SetRedis(ctx, cacheKey, dataCache, time.Minute*30)

	return userVouchers, nil
}
