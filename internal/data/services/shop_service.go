package services

import (
	"context"
	"nusantara_service/internal/domain/entities"
	dto "nusantara_service/internal/dto/request"

	"github.com/google/uuid"
)

type ShopService interface {
	CreateShop(ctx context.Context, superAdminId uuid.UUID, req dto.CreateShopRequest) error
	GetAllShop(ctx context.Context, page, limit int, search string) ([]*entities.ShopEntity, int, error)
	GetByIdShop(ctx context.Context, id uuid.UUID) (*entities.ShopEntity, error)
	UpdateShop(ctx context.Context, superAdminId uuid.UUID, id uuid.UUID, req dto.UpdateShopRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatusShop(ctx context.Context, id uuid.UUID, req dto.UpdateStatusShopRequest) error
}
