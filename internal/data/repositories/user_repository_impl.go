package repositories

import (
	"context"
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
func (u *UserRepositoryImpl) Create(ctx context.Context, user *entities.UserEntity) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}
