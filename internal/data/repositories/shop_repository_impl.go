package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopRepositoryImpl struct {
	db *gorm.DB
}

func NewShopRepositoryImpl(db *gorm.DB) repositories.ShopRepository {
	return &ShopRepositoryImpl{db: db}
}

func (s *ShopRepositoryImpl) preloadRelations(db *gorm.DB) *gorm.DB {
	return db.
		Preload("ShopImages.Image").
		Preload("ShopProducts.Product").
		Preload("ShopProducts.Product.Image").
		Preload("ShopProducts.Product.ProductImages.Image").
		Preload("ShopProducts.Product.TypeProduct").
		Preload("ShopProducts.Product.User").
		Preload("ShopProducts.Product.User.Role").
		Preload("ShopCashiers.User").
		Preload("ShopCashiers.User.Role").
		Preload("Creator").
		Preload("Creator.Role")
}

// Create implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) Create(ctx context.Context, data *entities.ShopEntity) error {
	return s.preloadRelations(s.db.WithContext(ctx)).Create(data).Error
}

// FindAll implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) FindAll(ctx context.Context, offset int, limit int, search string) ([]*entities.ShopEntity, int, error) {
	var (
		items []*entities.ShopEntity
		total int64
	)

	query := s.db.WithContext(ctx).Model(&entities.ShopEntity{}).Where("deleted_at IS NULL")
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := s.preloadRelations(query.Offset(offset).Limit(limit).Order("created_at DESC")).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, int(total), nil
}

// FindById implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*entities.ShopEntity, error) {
	var shop entities.ShopEntity
	if err := s.preloadRelations(s.db.WithContext(ctx)).First(&shop, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &shop, nil
}

// Update implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data *entities.ShopEntity) error {
	return s.preloadRelations(s.db.WithContext(ctx)).Save(data).Error
}

// Delete implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).Delete(&entities.ShopEntity{}, "id = ?", id).Error
}

// UpdateStatus implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status int) error {
	return s.preloadRelations(s.db.WithContext(ctx)).Model(&entities.ShopEntity{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// FindByName implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) FindByName(ctx context.Context, name string) (*entities.ShopEntity, error) {
	var shop entities.ShopEntity
	if err := s.preloadRelations(s.db.WithContext(ctx)).First(&shop, "name = ? AND deleted_at IS NULL", name).Error; err != nil {
		return nil, err
	}

	return &shop, nil
}

// FindUserIDSuperAdmin implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) FindUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := s.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND roles.name = ?", userID, "superadmin").
		Preload("Role").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindCashier implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) FindCashier(ctx context.Context, cashierID uuid.UUID) (*entities.UserEntity, error) {
	var cashier entities.UserEntity
	if err := s.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND roles.name = ?", cashierID, "admin").
		Preload("Role").
		First(&cashier).Error; err != nil {
		return nil, err
	}

	return &cashier, nil
}

// CreateGallery implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) CreateGallery(ctx context.Context, shopID uuid.UUID, imageID uuid.UUID, alt string) error {
	shopImage := &entities.ShopImageEntity{
		ID:      uuid.New(),
		ShopID:  shopID,
		ImageID: imageID,
		Altext:  alt,
	}
	return s.db.WithContext(ctx).Create(shopImage).Error
}

// CreateImage implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) CreateImage(ctx context.Context, img *entities.ImageEntity) error {
	return s.db.WithContext(ctx).Create(img).Error
}

// DeleteGalleryByShopID implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) DeleteGalleryByShopID(ctx context.Context, shopID uuid.UUID) error {
	return s.db.WithContext(ctx).Where("shop_id = ?", shopID).Delete(&entities.ShopImageEntity{}).Error
}

// UpdateShopCover implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) UpdateShopCover(ctx context.Context, shopID uuid.UUID, coverURL string) error {
	return s.db.WithContext(ctx).Model(&entities.ShopEntity{}).
		Where("id = ?", shopID).
		Update("cover", coverURL).Error
}

// AssignCashiers implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) AssignCashiers(ctx context.Context, shopID uuid.UUID, cashierIDs []uuid.UUID) error {
	if len(cashierIDs) == 0 {
		return nil
	}

	var rows []entities.ShopCashierEntity

	for _, id := range cashierIDs {
		rows = append(rows, entities.ShopCashierEntity{
			ID:         uuid.New(),
			ShopID:     shopID,
			CashierID:  id,
			AssignedAt: time.Now(),
		})
	}

	return s.db.WithContext(ctx).Create(&rows).Error
}

// AssignProducts implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) AssignProducts(ctx context.Context, shopID uuid.UUID, items []entities.ShopProductEntity) error {
	if len(items) == 0 {
		return nil
	}

	for i := range items {
		items[i].ID = uuid.New()
		items[i].ShopID = shopID
		items[i].AssignedAt = time.Now()
	}

	return s.db.WithContext(ctx).Create(&items).Error
}

// DeleteShopCashiersByShopID implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) DeleteShopCashiersByShopID(ctx context.Context, shopID uuid.UUID) error {
	return s.db.WithContext(ctx).Where("shop_id = ?", shopID).Delete(&entities.ShopCashierEntity{}).Error
}

// DeleteShopProductsByShopID implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) DeleteShopProductsByShopID(ctx context.Context, shopID uuid.UUID) error {
	return s.db.WithContext(ctx).Where("shop_id = ?", shopID).Delete(&entities.ShopProductEntity{}).Error
}

// ReplaceCashiers implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) ReplaceCashiers(ctx context.Context, shopID uuid.UUID, cashierIDs []uuid.UUID) error {
	if err := s.DeleteShopCashiersByShopID(ctx, shopID); err != nil {
		return err
	}

	if len(cashierIDs) == 0 {
		return nil
	}

	var rows []entities.ShopCashierEntity
	for _, id := range cashierIDs {
		rows = append(rows, entities.ShopCashierEntity{
			ID:         uuid.New(),
			ShopID:     shopID,
			CashierID:  id,
			AssignedAt: time.Now(),
		})
	}

	return s.db.WithContext(ctx).Create(&rows).Error
}

// ReplaceProducts implements repositories.ShopRepository.
func (s *ShopRepositoryImpl) ReplaceProducts(ctx context.Context, shopID uuid.UUID, items []entities.ShopProductEntity) error {
	if err := s.DeleteShopProductsByShopID(ctx, shopID); err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	for i := range items {
		items[i].ID = uuid.New()
		items[i].ShopID = shopID
		items[i].AssignedAt = time.Now()
	}

	return s.db.WithContext(ctx).Create(&items).Error
}
