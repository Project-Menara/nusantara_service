package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CashierRepositoryImpl struct {
	db *gorm.DB
}

func NewCashierRepositoryImpl(db *gorm.DB) repositories.CashierRepository {
	return &CashierRepositoryImpl{db: db}
}

// Create implements repositories.CashierRepository.
func (c *CashierRepositoryImpl) Create(ctx context.Context, data *entities.UserEntity) error {
	return c.db.WithContext(ctx).Create(data).Error
}

// FindByAll implements repositories.CashierRepository.
func (c *CashierRepositoryImpl) FindByAll(ctx context.Context, offset int, limit int, search string) ([]*entities.UserEntity, int, error) {
	var (
		cashier []*entities.UserEntity
		total   int64
	)

	query := c.db.WithContext(ctx).Model(&entities.UserEntity{}).Preload("Role").Where("deleted_at IS NULL").Joins("JOIN roles ON roles.id = users.role_id").Where("roles.name = ?", "admin")

	if search != "" {
		query = query.Where("users.name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("users.created_at DESC").Limit(limit).Offset(offset).Preload("Role").Find(&cashier).Error; err != nil {
		return nil, 0, err
	}

	return cashier, int(total), nil
}

// FindById implements repositories.CashierRepository.
func (c *CashierRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*entities.UserEntity, error) {
	var cashier entities.UserEntity
	if err := c.db.WithContext(ctx).Preload("Role").First(&cashier, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &cashier, nil
}

// Update implements repositories.CashierRepository.
func (c *CashierRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data *entities.UserEntity) error {
	return c.db.WithContext(ctx).Save(data).Error
}

// Delete implements repositories.CashierRepository.
func (c *CashierRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := c.db.WithContext(ctx).Delete(&entities.UserEntity{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// FindByEmail implements repositories.CashierRepository.
func (c *CashierRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.UserEntity, error) {
	var cashier entities.UserEntity
	if err := c.db.WithContext(ctx).Preload("Role").First(&cashier, "email = ? AND deleted_at IS NULL", email).Error; err != nil {
		return nil, err
	}

	return &cashier, nil
}

// FindByUsername implements repositories.CashierRepository.
func (c *CashierRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entities.UserEntity, error) {
	var cashier entities.UserEntity
	if err := c.db.WithContext(ctx).First(&cashier, "username = ? AND deleted_at IS NULL", username).Error; err != nil {
		return nil, err
	}
	return &cashier, nil
}
