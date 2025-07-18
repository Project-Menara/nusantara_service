package services

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/dto"

	"github.com/google/uuid"
)

type RoleService interface {
	CreateRole(ctx context.Context, req dto.CreateRoleRequest) (*entities.RoleEntity, error)
	GetAllRole(ctx context.Context) ([]*entities.RoleEntity, error)
	GetByIdRole(ctx context.Context, id uuid.UUID) (*entities.RoleEntity, error)
	UpdateRole(ctx context.Context, req dto.UpdateRoleRequest) (*entities.RoleEntity, error)
	DeleteRole(ctx context.Context, id uuid.UUID) error
}
