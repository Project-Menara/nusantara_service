package usecases

import (
	"context"
	"errors"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	dto "nusantara_service/internal/dto/request"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type roleService struct {
	repo repositories.RoleRepository
	rdb  *redis.Client
}

func NewRoleUsecase(repo repositories.RoleRepository, rdb *redis.Client) services.RoleService {
	return &roleService{repo: repo, rdb: rdb}
}

func (r *roleService) CreateRole(ctx context.Context, req dto.CreateRoleRequest) (*entities.RoleEntity, error) {
	existingRole, err := r.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingRole != nil {
		return nil, errors.New("role already exists")
	}

	newRole := &entities.RoleEntity{
		Name: req.Name,
	}

	_, err = r.repo.Create(ctx, newRole)
	if err != nil {
		return nil, errors.New("failed to add role")
	}

	return newRole, nil
}

// DeleteRole implements services.RoleService.
func (r *roleService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	_, err := r.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	return r.repo.Delete(ctx, id)

}

// GetAllRole implements services.RoleService.
func (r *roleService) GetAllRole(ctx context.Context, page, limit int) ([]*entities.RoleEntity, int, error) {
	offset := (page - 1) * limit
	roles, err := r.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := r.repo.CountAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// GetByIdRole implements services.RoleService.
func (r *roleService) GetByIdRole(ctx context.Context, id uuid.UUID) (*entities.RoleEntity, error) {
	role, err := r.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return role, nil
}

// UpdateRole implements services.RoleService.
func (r *roleService) UpdateRole(ctx context.Context, id uuid.UUID, req dto.UpdateRoleRequest) (*entities.RoleEntity, error) {
	role, err := r.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	existingRole, err := r.repo.FindByName(ctx, req.Name)
	if err == nil && existingRole.ID != id {
		return nil, errors.New("role already exists")
	}

	role.Name = req.Name
	updatedRole, err := r.repo.Update(ctx, id, role)
	if err != nil {
		return nil, err
	}

	return updatedRole, nil

}
