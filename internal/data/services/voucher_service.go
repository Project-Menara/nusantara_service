package services

import (
	"context"
	"nusantara_service/internal/domain/entities"
	dto "nusantara_service/internal/dto/request"

	"github.com/google/uuid"
)

type VoucherService interface {
	CreateVoucher(ctx context.Context, userId string, req dto.CreateVoucherRequest) (*entities.VoucherEntity, error)
	GetAllVoucher(ctx context.Context, page, limit int, search string) ([]*entities.VoucherEntity, int, error)
	GetByIdVoucher(ctx context.Context, id uuid.UUID) (*entities.VoucherEntity, error)
	UpdateVoucher(ctx context.Context, userId string, id uuid.UUID, req dto.UpdateVoucherRequest) (*entities.VoucherEntity, error)
	DeleteVoucher(ctx context.Context, id uuid.UUID) error
	UpdateStatusVoucher(ctx context.Context, id uuid.UUID, req dto.UpdateStatusVoucherRequest) error

	GetAllVoucherCustomer(ctx context.Context, page, limit int) ([]*entities.VoucherEntity, int, error)
	GetByIdVoucherCustomer(ctx context.Context, id uuid.UUID) (*entities.VoucherEntity, error)
}
