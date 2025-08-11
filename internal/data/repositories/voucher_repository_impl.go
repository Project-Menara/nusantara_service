package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VoucherRepositoryImpl struct {
	db *gorm.DB
}

func NewVoucherRepositoryImpl(db *gorm.DB) repositories.VoucherRepository {
	return &VoucherRepositoryImpl{db: db}
}

// CountAll implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) CountAll(ctx context.Context, search string) (int, error) {
	var count int64
	query := v.db.WithContext(ctx).Table("vouchers").Where("deleted_at IS NULL")
	if search != "" {
		query = query.Where("code ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// FindByCode implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) FindByCode(ctx context.Context, code string) (*entities.VoucherEntity, error) {
	var voucher entities.VoucherEntity
	if err := v.db.WithContext(ctx).First(&voucher, "code = ? AND deleted_at IS NULL", code).Error; err != nil {
		return nil, err
	}

	return &voucher, nil
}

// Create implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) Create(ctx context.Context, voucher *entities.VoucherEntity) (*entities.VoucherEntity, error) {
	err := v.db.WithContext(ctx).Create(voucher)
	if err != nil {
		return nil, err.Error
	}

	return voucher, nil
}

// FindAll implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) FindAll(ctx context.Context, limit int, offset int, search string) ([]*entities.VoucherEntity, error) {
	var voucher []*entities.VoucherEntity

	query := v.db.WithContext(ctx).Table("vouchers").Preload("User").Preload("User.Role").Where("deleted_at IS NULL").Order("created_at DESC").Limit(limit).Offset(offset)
	if search != "" {
		query = query.Where("code ILIKE ?", "%"+search+"%")
	}
	if err := query.Find(&voucher).Error; err != nil {
		return nil, err
	}

	return voucher, nil
}

// FindById implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*entities.VoucherEntity, error) {
	var voucher entities.VoucherEntity
	if err := v.db.WithContext(ctx).Preload("User").Preload("User.Role").First(&voucher, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &voucher, nil
}

// Update implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data *entities.VoucherEntity) (*entities.VoucherEntity, error) {
	var voucher entities.VoucherEntity
	if err := v.db.WithContext(ctx).Preload("User").Preload("User.Role").First(&voucher, "id = ?", id).Error; err != nil {
		return nil, err
	}

	voucher.Code = data.Code
	voucher.DiscountAmount = data.DiscountAmount
	voucher.DiscountPercent = data.DiscountPercent
	voucher.MinimumSpend = data.MinimumSpend
	voucher.StartDate = data.StartDate
	voucher.EndDate = data.EndDate
	voucher.Quota = data.Quota
	voucher.Description = data.Description
	voucher.DiscountType = data.DiscountType

	if err := v.db.WithContext(ctx).Updates(&voucher).Error; err != nil {
		return nil, err
	}

	return &voucher, nil
}

// Delete implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := v.db.WithContext(ctx).Delete(&entities.VoucherEntity{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// FindByUserIDSuperAdmin implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) FindByUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := v.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND roles.name = ?", userID, "superadmin").
		Preload("Role").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateStatus implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status int) error {
	return v.db.WithContext(ctx).Model(&entities.VoucherEntity{}).Where("id = ?", id).Update("status", status).Error
}

// GetAllVoucherCustomer implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) GetAllVoucherCustomer(ctx context.Context, limit int, offset int) ([]*entities.VoucherEntity, error) {
	var voucher []*entities.VoucherEntity

	if err := v.db.WithContext(ctx).Table("vouchers").Preload("User").Preload("User.Role").Where("deleted_at IS NULL AND status = 1").Order("created_at DESC").Limit(limit).Offset(offset).Find(&voucher).Error; err != nil {
		return nil, err
	}

	return voucher, nil
}

// GetByIdVoucherCustomer implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) GetByIdVoucherCustomer(ctx context.Context, id uuid.UUID) (*entities.VoucherEntity, error) {
	var voucher entities.VoucherEntity
	if err := v.db.WithContext(ctx).Preload("User").Preload("User.Role").Where("status = 1").First(&voucher, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &voucher, nil
}

// CheckUserVoucherClaimed implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) CheckUserVoucherClaimed(ctx context.Context, customerID uuid.UUID, voucherID uuid.UUID) (bool, error) {
	var count int64
	err := v.db.WithContext(ctx).Model(&entities.UserVoucherEntity{}).
		Where("user_id = ? AND voucher_id = ?", customerID, voucherID).
		Count(&count).Error

	return count > 0, err
}

// IncreaseVoucherClaimedCount implements repositories.VoucherRepository.
func (v *VoucherRepositoryImpl) IncreaseVoucherClaimedCount(ctx context.Context, voucherID uuid.UUID) error {
	return v.db.WithContext(ctx).Model(&entities.VoucherEntity{}).
		Where("id = ?", voucherID).
		Update("claimed_count", gorm.Expr("claimed_count + ?", 1)).Error
}
