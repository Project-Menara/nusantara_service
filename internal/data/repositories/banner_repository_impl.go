package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BannerRepositoryImpl struct {
	db *gorm.DB
}

func NewBannerRepositoryImpl(db *gorm.DB) repositories.BannerRepository {
	return &BannerRepositoryImpl{db: db}
}

// CountAll implements repositories.BannerRepository.
func (b *BannerRepositoryImpl) CountAll(ctx context.Context) (int, error) {
	var count int64
	var banners entities.BannerEntity
	if err := b.db.WithContext(ctx).Model(&banners).Count(&count).Error; err != nil {
		return 0, nil
	}

	return int(count), nil
}

// FindByName implements repositories.BannerRepository.
func (b *BannerRepositoryImpl) FindByName(ctx context.Context, name string) (*entities.BannerEntity, error) {
	var banner entities.BannerEntity
	if err := b.db.WithContext(ctx).First(&banner, "name = ? AND deleted_at IS NULL", name).Error; err != nil {
		return nil, err
	}

	return &banner, nil
}

// Create implements repositories.BannerRepository.
func (b *BannerRepositoryImpl) Create(ctx context.Context, banner *entities.BannerEntity) (*entities.BannerEntity, error) {
	err := b.db.WithContext(ctx).Create(banner)
	if err != nil {
		return nil, err.Error
	}

	return banner, nil
}

// FindAll implements repositories.BannerRepository.
func (b *BannerRepositoryImpl) FindAll(ctx context.Context, limit int, offset int) ([]*entities.BannerEntity, error) {
	var banners []*entities.BannerEntity
	err := b.db.WithContext(ctx).
		Order("updated_at DESC"). // Urut dari yang terbaru
		Preload("User").
		Preload("User.Role").
		Limit(limit).
		Offset(offset).
		Find(&banners).Error

	if err != nil {
		return nil, err
	}
	return banners, nil
}

// FindById implements repositories.BannerRepository.
func (b *BannerRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*entities.BannerEntity, error) {
	var banner entities.BannerEntity
	if err := b.db.WithContext(ctx).Preload("User").Preload("User.Role").First(&banner, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &banner, nil
}

// Update implements repositories.BannerRepository.
func (b *BannerRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data *entities.BannerEntity) (*entities.BannerEntity, error) {
	var banner entities.BannerEntity
	if err := b.db.WithContext(ctx).Preload("User").Preload("User.Role").First(&banner, "id = ?", id).Error; err != nil {
		return nil, err
	}

	banner.Name = data.Name
	banner.Description = data.Description
	banner.Photo = data.Photo
	banner.UserID = data.UserID

	if err := b.db.WithContext(ctx).Updates(&banner).Error; err != nil {
		return nil, err
	}

	return &banner, nil
}

// Delete implements repositories.BannerRepository.
func (b *BannerRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := b.db.WithContext(ctx).Delete(&entities.BannerEntity{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// FindByUserIDSuperAdmin implements repositories.UserRepository.
func (u *BannerRepositoryImpl) FindByUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND roles.name = ?", userID, "superadmin").
		Preload("Role").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil

}

// UpdateStatus implements repositories.BannerRepository.
func (b *BannerRepositoryImpl) UpdateStatus(ctx context.Context, bannerId string, status int) error {
	return b.db.WithContext(ctx).Model(&entities.BannerEntity{}).
		Where("id = ?", bannerId).
		Update("status", status).Error
}

// GetAllCustomer implements repositories.BannerRepository.
func (b *BannerRepositoryImpl) GetAllBannerCustomer(ctx context.Context) ([]*entities.BannerEntity, error) {
	var banners []*entities.BannerEntity
	if err := b.db.WithContext(ctx).Preload("User").Preload("User.Role").Where("status = 1").Find(&banners).Error; err != nil {
		return nil, err
	}
	return banners, nil
}
