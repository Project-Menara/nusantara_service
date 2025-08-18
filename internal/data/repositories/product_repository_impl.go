package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepositoryImpl(db *gorm.DB) repositories.ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

func (p *ProductRepositoryImpl) preloadRelations(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Image").
		Preload("ProductImages.Image").
		Preload("TypeProduct").
		Preload("TypeProduct.User").
		Preload("TypeProduct.User.Role").
		Preload("User").
		Preload("User.Role")
}

// Create implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) Create(ctx context.Context, product *entities.ProductEntity) error {
	return p.preloadRelations(p.db.WithContext(ctx)).Create(product).Error
}

// GetByID implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.ProductEntity, error) {
	var product entities.ProductEntity
	err := p.preloadRelations(p.db.WithContext(ctx)).First(&product, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetAll implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) GetAll(ctx context.Context, offset int, limit int, search string) ([]*entities.ProductEntity, int, error) {
	var (
		items []*entities.ProductEntity
		total int64
	)
	q := p.db.WithContext(ctx).Model(&entities.ProductEntity{}).Where("deleted_at IS NULL")
	if search != "" {
		q = q.Where("name ILIKE ?", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := p.preloadRelations(
		q.Offset(offset).Limit(limit).Order("created_at DESC"),
	).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, int(total), nil
}

// Update implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) Update(ctx context.Context, id uuid.UUID, product *entities.ProductEntity) error {
	return p.preloadRelations(p.db.WithContext(ctx)).Save(product).Error
}

// Delete implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return p.db.WithContext(ctx).Delete(&entities.ProductEntity{}, "id = ?", id).Error
}

// CreateImages implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) CreateImages(ctx context.Context, images []entities.ImageEntity) error {
	return p.db.WithContext(ctx).Create(&images).Error
}

// CreateProductImages implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) CreateProductImages(ctx context.Context, piv []entities.ProductImageEntity) error {
	return p.db.WithContext(ctx).Create(&piv).Error
}

// DeleteProductImages implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) DeleteProductImages(ctx context.Context, productID uuid.UUID) error {
	return p.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&entities.ProductImageEntity{}).Error
}

// GetProductImagePublicIDs implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) GetProductImagePublicIDs(ctx context.Context, productID uuid.UUID) ([]string, error) {
	var piv []*entities.ProductImageEntity
	if err := p.db.WithContext(ctx).Preload("Image").Where("product_id = ?", productID).Find(&piv).Error; err != nil {
		return nil, err
	}

	out := make([]string, 0, len(piv))
	for _, pi := range piv {
		out = append(out, pi.Image.ImagePath)
	}
	return out, nil
}

// FindByUserIDSuperAdmin implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) FindByUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := p.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND roles.name = ?", userID, "superadmin").
		Preload("Role").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByName implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) FindByName(ctx context.Context, name string) (*entities.ProductEntity, error) {
	var product entities.ProductEntity
	if err := p.preloadRelations(p.db.WithContext(ctx)).First(&product, "name = ? AND deleted_at IS NULL", name).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// UpdateStatus implements repositories.ProductRepository.
func (p *ProductRepositoryImpl) UpdateStatus(ctx context.Context, productID uuid.UUID, status int) error {
	return p.preloadRelations(p.db.WithContext(ctx)).Model(&entities.ProductEntity{}).
		Where("id = ?", productID).
		Update("status", status).Error
}
