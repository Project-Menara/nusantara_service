package repositories

import (
	"context"
	"log"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// SUPERADMIN
// CreateAdmin implements repositories.UserRepository.
func (u *UserRepositoryImpl) CreateAdmin(ctx context.Context, user *entities.UserEntity) (*entities.UserEntity, error) {
	err := u.db.WithContext(ctx).Create(user).Error // Chain .Error here
	if err != nil {
		return nil, err // Corrected: return the error itself
	}

	// Make sure the user variable is updated by passing its address
	err = u.db.WithContext(ctx).Preload("Role").First(user, "id = ?", user.ID).Error // Chain .Error and pass user directly
	if err != nil {
		log.Println("Failed to preload role:", err) // Changed "Gagal" to "Failed" for consistency
		return nil, err                             // Corrected: return the error itself
	}

	return user, nil
}

// FindExistUsername implements repositories.UserRepository.
func (u *UserRepositoryImpl) FindExistUsername(ctx context.Context, username string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail implements repositories.UserRepository.
func (u *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).Preload("Role").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindUserById implements repositories.UserRepository.
func (u *UserRepositoryImpl) FindUserById(ctx context.Context, userId string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).Preload("Role").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ChangePassword implements repositories.UserRepository.
func (u *UserRepositoryImpl) ChangePassword(ctx context.Context, userId string, Updated *entities.UserEntity) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).Preload("Role").First(&user, "id = ?", userId).Where("role").Error; err != nil {
		return nil, err
	}

	user.Password = Updated.Password

	if err := u.db.WithContext(ctx).Updates(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

//END SUPERADMIN
