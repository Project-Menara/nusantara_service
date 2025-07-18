package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepositoryImpl(db *gorm.DB) repositories.RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

func (r *RoleRepositoryImpl) Create(ctx context.Context, role *entities.RoleEntity) (*entities.RoleEntity, error) {
	err := r.db.WithContext(ctx).Create(role)
	if err != nil {
		return nil, err.Error
	}

	return role, nil
}

// Delete implements repositories.RoleRepository.
func (r *RoleRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&entities.RoleEntity{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// FindAll implements repositories.RoleRepository.
func (r *RoleRepositoryImpl) FindAll(ctx context.Context) ([]*entities.RoleEntity, error) {
	var roles []*entities.RoleEntity
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

// FindById implements repositories.RoleRepository.
func (r *RoleRepositoryImpl) FindById(ctx context.Context, Id uuid.UUID) (*entities.RoleEntity, error) {
	var role entities.RoleEntity
	if err := r.db.WithContext(ctx).First(&role, "id = ?", Id).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// FindByName implements repositories.RoleRepository.
func (r *RoleRepositoryImpl) FindByName(ctx context.Context, name string) (*entities.RoleEntity, error) {
	var role entities.RoleEntity
	if err := r.db.WithContext(ctx).First(&role, "name = ?", name).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// Update implements repositories.RoleRepository.
func (r *RoleRepositoryImpl) Update(ctx context.Context, id uuid.UUID, updated *entities.RoleEntity) (*entities.RoleEntity, error) {
	var role entities.RoleEntity
	if err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}

	role.Name = updated.Name

	if err := r.db.WithContext(ctx).Updates(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
