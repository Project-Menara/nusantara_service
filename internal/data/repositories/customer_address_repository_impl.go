package repositories

import (
	"context"
	"errors"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerAddressRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerAddresRepositoryImpl(db *gorm.DB) repositories.CustomerAddressRepository {
	return &CustomerAddressRepositoryImpl{db: db}
}

func (c *CustomerAddressRepositoryImpl) preloadRelations(db *gorm.DB) *gorm.DB {
	return db.
		Preload("User").
		Preload("User.Role")
}

func (c *CustomerAddressRepositoryImpl) preloadRelationsShop(db *gorm.DB) *gorm.DB {
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

// Create implements repositories.CustomerAddressRepository.
func (c *CustomerAddressRepositoryImpl) Create(ctx context.Context, data *entities.CustomerAddressEntity) error {
	return c.preloadRelations(c.db.WithContext(ctx)).Create(data).Error
}

// FindByID implements repositories.CustomerAddressRepository.
func (c *CustomerAddressRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.CustomerAddressEntity, error) {
	var addr entities.CustomerAddressEntity
	if err := c.preloadRelations(c.db.WithContext(ctx)).First(&addr, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}

	return &addr, nil
}

// FindByUser implements repositories.CustomerAddressRepository.
func (c *CustomerAddressRepositoryImpl) FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.CustomerAddressEntity, error) {
	var list []*entities.CustomerAddressEntity

	if err := c.preloadRelations(c.db.WithContext(ctx)).Where("user_id = ?", userID).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

// Update implements repositories.CustomerAddressRepository.
func (c *CustomerAddressRepositoryImpl) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, data *entities.CustomerAddressEntity) error {
	return c.preloadRelations(c.db.WithContext(ctx)).Model(&entities.CustomerAddressEntity{}).Where("id = ? AND user_id = ?", id, userID).Updates(data).Error
}

// Delete implements repositories.CustomerAddressRepository.
func (c *CustomerAddressRepositoryImpl) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return c.db.WithContext(ctx).Delete(&entities.CustomerAddressEntity{}, "id = ? AND user_id = ?", id, userID).Error
}

// FindCustomerID implements repositories.CustomerAddressRepository.
func (c *CustomerAddressRepositoryImpl) FindCustomerID(ctx context.Context, userID string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := c.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND roles.name = ?", userID, "customer").
		Preload("Role").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ClearDefaultByUser implements repositories.CustomerAddressRepository.
func (c *CustomerAddressRepositoryImpl) ClearDefaultByUser(ctx context.Context, userID uuid.UUID) error {
	return c.db.WithContext(ctx).
		Model(&entities.CustomerAddressEntity{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}

// FindDefaultByUser implements repositories.CustomerAddressRepository.
func (c *CustomerAddressRepositoryImpl) FindDefaultByUser(ctx context.Context, userID uuid.UUID) (*entities.CustomerAddressEntity, error) {
	var addr entities.CustomerAddressEntity
	err := c.preloadRelations(c.db.WithContext(ctx)).
		Where("user_id = ? AND is_default = true", userID).
		First(&addr).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // biar jelas, ga dianggap error
	}
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

// SetDefaultByUser implements repositories.CustomerAddressRepository.
func (c *CustomerAddressRepositoryImpl) SetDefaultByUser(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error {
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entities.CustomerAddressEntity{}).
			Where("user_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return err
		}
		// pastikan address memang milik user (WHERE id = ? AND user_id = ?)
		if err := tx.Model(&entities.CustomerAddressEntity{}).
			Where("id = ? AND user_id = ?", addressID, userID).
			Update("is_default", true).Error; err != nil {
			return err
		}
		return nil
	})
}
func (c *CustomerAddressRepositoryImpl) SetFirstAsDefault(ctx context.Context, userID uuid.UUID) error {
	var addr entities.CustomerAddressEntity
	err := c.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&addr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// tidak ada alamat, nothing to set
			return nil
		}
		return err
	}
	// gunakan SetDefaultByUser (transactional)
	return c.SetDefaultByUser(ctx, userID, addr.ID)
}

// FindNearby implements repositories.CustomerAddressRepository.
func (c *CustomerAddressRepositoryImpl) FindNearby(ctx context.Context, lat float64, lng float64, radiusKm float64) ([]*entities.ShopEntity, map[uuid.UUID]float64, error) {
	type idDist struct {
		ID       uuid.UUID `gorm:"column:id"`
		Distance float64   `gorm:"column:distance"`
	}

	var results []idDist
	query := `
		SELECT id::uuid AS id,
			(6371 * acos(
				cos(radians(?)) * cos(radians(lat)) * cos(radians(lng) - radians(?)) +
				sin(radians(?)) * sin(radians(lat))
			)) AS distance
		FROM shops
		WHERE status = 1
		AND (6371 * acos(
				cos(radians(?)) * cos(radians(lat)) * cos(radians(lng) - radians(?)) +
				sin(radians(?)) * sin(radians(lat))
			)) <= ?
		ORDER BY distance ASC
		LIMIT 20;
	`

	if err := c.db.WithContext(ctx).Raw(query,
		lat, lng, lat, // untuk SELECT
		lat, lng, lat, // untuk AND
		radiusKm,
	).Scan(&results).Error; err != nil {
		return nil, nil, err
	}

	if len(results) == 0 {
		return []*entities.ShopEntity{}, map[uuid.UUID]float64{}, nil
	}

	ids := make([]uuid.UUID, 0, len(results))
	orderMap := make(map[uuid.UUID]float64, len(results))
	for _, r := range results {
		ids = append(ids, r.ID)
		orderMap[r.ID] = r.Distance
	}

	// Ambil shops dengan preload relations
	var shops []*entities.ShopEntity
	if err := c.preloadRelationsShop(c.db.WithContext(ctx)).
		Where("id IN ?", ids).
		Find(&shops).Error; err != nil {
		return nil, nil, err
	}

	// Reorder sesuai urutan distance
	ordered := make([]*entities.ShopEntity, 0, len(shops))
	shopMap := make(map[uuid.UUID]*entities.ShopEntity, len(shops))
	for _, s := range shops {
		shopMap[s.ID] = s
	}
	for _, id := range ids {
		if sh, ok := shopMap[id]; ok {
			ordered = append(ordered, sh)
		}
	}

	return ordered, orderMap, nil
}
