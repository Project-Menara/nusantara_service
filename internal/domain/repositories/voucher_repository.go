package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"

	"github.com/google/uuid"
)

type VoucherRepository interface {
	Create(ctx context.Context, voucher *entities.VoucherEntity) (*entities.VoucherEntity, error)
	FindAll(ctx context.Context, limit, offset int, search string) ([]*entities.VoucherEntity, error)
	FindById(ctx context.Context, id uuid.UUID) (*entities.VoucherEntity, error)
	Update(ctx context.Context, id uuid.UUID, data *entities.VoucherEntity) (*entities.VoucherEntity, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status int) error

	FindByCode(ctx context.Context, code string) (*entities.VoucherEntity, error)
	CountAll(ctx context.Context, search string) (int, error)
	FindByUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error)

	GetAllVoucherCustomer(ctx context.Context, limit, offset int) ([]*entities.VoucherEntity, error)
	GetByIdVoucherCustomer(ctx context.Context, id uuid.UUID) (*entities.VoucherEntity, error)
}
