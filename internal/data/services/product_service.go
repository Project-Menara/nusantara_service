package services

import (
	"context"
	"mime/multipart"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/dto"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, userId uuid.UUID, req dto.CreateProductRequest) (*entities.ProductEntity, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*entities.ProductEntity, error)
	GetProductAll(ctx context.Context, page, limit int, search string) ([]*entities.ProductEntity, int, error)
	UpdateProduct(ctx context.Context, userId uuid.UUID, req dto.UpdateProductRequest) (*entities.ProductEntity, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UploadMany(ctx context.Context, folder, prefix string, files []*multipart.FileHeader, workers int) ([]string, []string, error)
	UpdateStatusProduct(ctx context.Context, productID uuid.UUID, req dto.UpdatStatusProducRequest) error
}
