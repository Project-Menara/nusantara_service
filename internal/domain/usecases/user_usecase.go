package usecases

import (
	"context"
	"errors"
	"fmt"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userService struct {
	repo repositories.UserRepository
	rdb  *redis.Client
}

func NewUserUsecase(repo repositories.UserRepository, rdb *redis.Client) services.UserService {
	return &userService{repo: repo, rdb: rdb}
}

//SUPERADMIN

// RegisterAdmin implements services.UserService.
func (u *userService) RegisterAdmin(ctx context.Context, req dto.RegisterAdminRequest) (*entities.UserEntity, error) {
	existingUsername, err := u.repo.FindExistUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingUsername != nil {
		return nil, errors.New("username already exists")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newAdmin := &entities.UserEntity{
		ID:       uuid.NewString(),
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
		RoleID:   uuid.MustParse(req.RoleID),
		Status:   1,
	}

	createdUser, err := u.repo.CreateAdmin(ctx, newAdmin)
	if err != nil {
		return nil, errors.New("failed to create admin")
	}

	return createdUser, nil

}

// LoginAdmin implements services.UserService.
func (u *userService) LoginAdmin(ctx context.Context, req dto.LoginAdminRequest) (string, error) {
	loginKey := fmt.Sprintf("login_attempts:%s", req.Email)

	const maxAttempts = 5
	const cooldownDuration = 1 * time.Minute

	attempts, err := u.rdb.Get(ctx, loginKey).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("failed to get login attempts from redis: %w", err)
	}

	ttl, err := u.rdb.TTL(ctx, loginKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("failed to get TTL from redis: %w", err)
	}

	if ttl > 0 && attempts >= maxAttempts {
		return "", &response.CooldownError{
			Message:    "Too many login attempts, please try again later.",
			RetryAfter: ttl,
		}
	}

	admin, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		u.rdb.Incr(ctx, loginKey)
		u.rdb.Expire(ctx, loginKey, cooldownDuration)
		return "", errors.New("email incorrect")
	}

	isPassword := utils.CheckPasswordHash(req.Password, admin.Password)
	if !isPassword {
		u.rdb.Incr(ctx, loginKey)
		u.rdb.Expire(ctx, loginKey, cooldownDuration)
		return "", errors.New("password incorrect")
	}

	u.rdb.Del(ctx, loginKey)

	token, err := utils.GenerateToken(admin.ID)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	expiry, err := utils.GetExpiryFromToken(token)
	if err != nil {
		return "", errors.New("failed to get token expiry")
	}

	// Menyimpan token ke Redis dengan waktu kedaluwarsa yang sama dengan token JWT
	redisKey := fmt.Sprintf("admin_token:%s", admin.ID)
	err = u.rdb.Set(ctx, redisKey, token, time.Until(expiry)).Err()
	if err != nil {
		return "", errors.New("failed to store token in redis")
	}

	return token, nil
}

func (u *userService) GetProfile(ctx context.Context, userId string, token string) (*entities.UserEntity, error) {
	redis_key := fmt.Sprintf("admin_token:%s", userId)
	storedToken, err := u.rdb.Get(ctx, redis_key).Result()
	if err != nil || storedToken != token {
		return nil, errors.New("invalid or expired session")
	}

	admin, err := u.repo.FindUserById(ctx, userId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return admin, nil
}

func (u *userService) LogoutAdmin(ctx context.Context, token string) error {
	expiry, err := utils.GetExpiryFromToken(token)
	if err != nil {
		return err
	}

	blackListKey := fmt.Sprintf("blacklist:%s", token)
	err = u.rdb.Set(ctx, blackListKey, "blacklisted", time.Until(expiry)).Err()
	if err != nil {
		return errors.New("failed to blacklist token")
	}

	return nil
}

// ChangePasswordSuperAdmin implements services.UserService.
func (u *userService) ChangePasswordSuperAdmin(ctx context.Context, userId string, token string, req dto.ChangePasswordRequest) (*entities.UserEntity, error) {
	redis_key := fmt.Sprintf("admin_token:%s", userId)
	storedToken, err := u.rdb.Get(ctx, redis_key).Result()
	if err != nil || storedToken != token {
		return nil, errors.New("invalid or expired session")
	}

	user, err := u.repo.FindUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return nil, errors.New("invalid current password")
	}
	if req.NewPassword != req.ConfirmationPassword {
		return nil, errors.New("new password and confirmation do not match")
	}
	if len(req.NewPassword) < 6 {
		return nil, errors.New("new password must be at least 6 characters long")
	}

	hashedNewPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, errors.New("failed to hash new password")
	}

	user.Password = hashedNewPassword

	updateData, err := u.repo.ChangePassword(ctx, userId, user)
	if err != nil {
		return nil, err
	}

	return updateData, nil

}

// CheckTokenBlacklisted implements services.UserService.
func (u *userService) CheckTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	val, err := u.rdb.Get(ctx, fmt.Sprintf("blacklist:%s", token)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil // Key does not exist, so not blacklisted
		}
		return false, err // Other actual Redis error
	}
	return val == "blacklisted", nil // Key exists, check its value

}
