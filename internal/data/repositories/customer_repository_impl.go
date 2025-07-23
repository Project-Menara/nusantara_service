package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"

	"gorm.io/gorm"
)

type CustomerRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerRepositoryImpl(db *gorm.DB) repositories.CustomerRepository {
	return &CustomerRepositoryImpl{db: db}
}

// CheckPhone implements repositories.UserRepository.
func (u *CustomerRepositoryImpl) FindByPhoneCustomer(ctx context.Context, phone string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.phone = ? AND roles.name = ?", phone, "customer").
		Preload("Role").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByName implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByPhone implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindByPhone(ctx context.Context, phone string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).First(&user, "phone = ?", phone).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindRoleByName implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindRoleByName(ctx context.Context, role string) (*entities.RoleEntity, error) {
	var roles entities.RoleEntity
	if err := u.db.WithContext(ctx).First(&roles, "name = ?", role).Error; err != nil {
		return nil, err
	}

	return &roles, nil
}

// RegisterCustomer implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) CreateCustomer(ctx context.Context, user *entities.UserEntity) (*entities.UserEntity, error) {
	err := u.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	err = u.db.WithContext(ctx).Preload("Role").First(user, "id = ?", user.ID).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
