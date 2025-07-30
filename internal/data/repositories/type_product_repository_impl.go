package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TypeProductRepositoryImpl struct {
	db *gorm.DB
}

func NewTypeProductRepositoryImpl(db *gorm.DB) repositories.TypeProductRepository {
	return &TypeProductRepositoryImpl{db: db}
}

// CountAll implements repositories.TypeProductRepository.
func (t *TypeProductRepositoryImpl) CountAll(ctx context.Context) (int, error) {
	var count int64
	var typeProducts entities.TypeProductEntity
	if err := t.db.WithContext(ctx).Model(&typeProducts).Count(&count).Error; err != nil {
		return 0, nil
	}

	return int(count), nil
}

// FindByName implements repositories.TypeProductRepository.
func (t *TypeProductRepositoryImpl) FindByName(ctx context.Context, name string) (*entities.TypeProductEntity, error) {
	var typeProduct entities.TypeProductEntity
	if err := t.db.WithContext(ctx).First(&typeProduct, "name = ? AND deleted_at IS NULL", name).Error; err != nil {
		return nil, err
	}

	return &typeProduct, nil
}

// Create implements repositories.TypeProductRepository.
func (t *TypeProductRepositoryImpl) Create(ctx context.Context, typeProduct *entities.TypeProductEntity) (*entities.TypeProductEntity, error) {
	err := t.db.WithContext(ctx).Create(typeProduct)
	if err != nil {
		return nil, err.Error
	}

	return typeProduct, nil
}

// FindAll implements repositories.TypeProductRepository.
func (t *TypeProductRepositoryImpl) FindAll(ctx context.Context, limit int, offset int) ([]*entities.TypeProductEntity, error) {
	var typeProducts []*entities.TypeProductEntity
	if err := t.db.WithContext(ctx).Preload("User").Preload("User.Role").Limit(limit).Offset(offset).Find(&typeProducts).Error; err != nil {
		return nil, err
	}

	return typeProducts, nil
}

// FindById implements repositories.TypeProductRepository.
func (t *TypeProductRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*entities.TypeProductEntity, error) {
	var typeProduct entities.TypeProductEntity
	if err := t.db.WithContext(ctx).Preload("User").Preload("User.Role").First(&typeProduct, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &typeProduct, nil
}

// Update implements repositories.TypeProductRepository.
func (t *TypeProductRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data *entities.TypeProductEntity) (*entities.TypeProductEntity, error) {
	var typeProduct entities.TypeProductEntity
	if err := t.db.WithContext(ctx).Preload("User").Preload("User.Role").First(&typeProduct, "id = ?", id).Error; err != nil {
		return nil, err
	}

	typeProduct.Name = data.Name
	typeProduct.Image = data.Image
	typeProduct.UserID = data.UserID

	if err := t.db.WithContext(ctx).Updates(&typeProduct).Error; err != nil {
		return nil, err
	}

	return &typeProduct, nil
}

// Delete implements repositories.TypeProductRepository.
func (t *TypeProductRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := t.db.WithContext(ctx).Delete(&entities.TypeProductEntity{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// FindByUserIDSuperAdmin implements repositories.TypeProductRepository.
func (t *TypeProductRepositoryImpl) FindByUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := t.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND roles.name = ?", userID, "superadmin").
		Preload("Role").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateStatus implements repositories.TypeProductRepository.
func (t *TypeProductRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status int) error {
	return t.db.WithContext(ctx).Model(&entities.TypeProductEntity{}).Where("id = ?", id).Update("status", status).Error
}

// GetAllTypeProductCustomer implements repositories.TypeProductRepository.
func (t *TypeProductRepositoryImpl) GetAllTypeProductCustomer(ctx context.Context) ([]*entities.TypeProductEntity, error) {
	var typeProducts []*entities.TypeProductEntity
	if err := t.db.WithContext(ctx).Preload("User").Preload("User.Role").Where("status = 1").Find(&typeProducts).Error; err != nil {
		return nil, err
	}

	return typeProducts, nil
}

// GetByIdTypeProductCustomer implements repositories.TypeProductRepository.
func (t *TypeProductRepositoryImpl) GetByIdTypeProductCustomer(ctx context.Context, id uuid.UUID) (*entities.TypeProductEntity, error) {
	var typeProduct entities.TypeProductEntity
	if err := t.db.WithContext(ctx).Preload("User").Preload("User.Role").Where("status = 1").First(&typeProduct, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &typeProduct, nil
}
