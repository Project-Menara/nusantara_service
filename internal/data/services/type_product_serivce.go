package services

import (
	"context"
	"mime/multipart"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/dto"

	"github.com/google/uuid"
)

type TypeProductService interface {
	CreateTypeProduct(ctx context.Context, userId string, req dto.CreateTypeProductRequest, image *multipart.FileHeader) (*entities.TypeProductEntity, error)
	GetAllTypeProduct(ctx context.Context, page, limit int, search string) ([]*entities.TypeProductEntity, int, error)
	GetByIdTypeProduct(ctx context.Context, id uuid.UUID) (*entities.TypeProductEntity, error)
	UpdateTypeProduct(ctx context.Context, userId string, id uuid.UUID, req dto.UpdateTypeProductRequest, image *multipart.FileHeader) (*entities.TypeProductEntity, error)
	DeleteTypeProduct(ctx context.Context, id uuid.UUID) error
	UpdateStatusTypeProduct(ctx context.Context, id uuid.UUID, req dto.UpdateStatusTypeProductRequest) error

	GetAllTypeProductCustomer(ctx context.Context) ([]*entities.TypeProductEntity, error)
	GetByIdTypeProductCustomer(ctx context.Context, id uuid.UUID) (*entities.TypeProductEntity, error)
}
