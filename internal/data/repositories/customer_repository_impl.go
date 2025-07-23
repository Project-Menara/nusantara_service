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
func (u *CustomerRepositoryImpl) FindByPhone(ctx context.Context, phone string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).First(&user, "phone = ?", phone).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
