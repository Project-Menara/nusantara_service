package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
)

type CustomerRepository interface {
	FindByPhone(ctx context.Context, phone string) (*entities.UserEntity, error)
}
