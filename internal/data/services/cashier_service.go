package services

import (
	"context"
	"nusantara_service/internal/domain/entities"
	dto "nusantara_service/internal/dto/request"

	"github.com/google/uuid"
)

type CashierService interface {
	CreateCashier(ctx context.Context, req dto.CreateCashierRequest) error
	GetCashierAll(ctx context.Context, page, limit int, search string) ([]*entities.UserEntity, int, error)
	GetCashierById(ctx context.Context, id uuid.UUID) (*entities.UserEntity, error)
	UpdateCashier(ctx context.Context, id uuid.UUID, req dto.UpdateCashierRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}
