package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerRepositoryImpl(db *gorm.DB) repositories.CustomerRepository {
	return &CustomerRepositoryImpl{db: db}
}

// CheckPhone implements repositories.UserRepository.
func (u *CustomerRepositoryImpl) FindByPhoneCustomer(ctx context.Context, phone string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.phone = ? AND roles.name = ?", phone, "customer").
		Preload("Role").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByName implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByPhone implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindByPhone(ctx context.Context, phone string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).First(&user, "phone = ?", phone).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindRoleByName implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindRoleByName(ctx context.Context, role string) (*entities.RoleEntity, error) {
	var roles entities.RoleEntity
	if err := u.db.WithContext(ctx).First(&roles, "name = ?", role).Error; err != nil {
		return nil, err
	}

	return &roles, nil
}

// RegisterCustomer implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) CreateCustomer(ctx context.Context, user *entities.UserEntity) (*entities.UserEntity, error) {
	err := u.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	err = u.db.WithContext(ctx).Preload("Role").First(user, "id = ?", user.ID).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateStatusCustomer implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) UpdateStatusCustomer(ctx context.Context, userID string, status int) error {
	return u.db.WithContext(ctx).Model(&entities.UserEntity{}).
		Where("id = ?", userID).
		Update("status", status).Error
}

// UpdatePinCustomer implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) UpdatePinCustomer(ctx context.Context, userID string, pin string) error {
	if err := u.db.WithContext(ctx).Model(&entities.UserEntity{}).Where("id = ?", userID).Update("password", pin).Error; err != nil {
		return err
	}
	return nil
}

// FindByUseIDCustomer implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindByUseIDCustomer(ctx context.Context, userID string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND roles.name = ?", userID, "customer").
		Preload("Role").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateCustomer implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) UpdateCustomer(ctx context.Context, userId string, data *entities.UserEntity) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).Preload("Role").First(&user, "id = ?", userId).Where("role").Error; err != nil {
		return nil, err
	}

	user.Name = data.Name
	user.Email = data.Email
	user.Username = data.Username
	user.Gender = data.Gender
	user.DateOfBirth = data.DateOfBirth
	user.Photo = data.Photo

	if err := u.db.WithContext(ctx).Updates(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ChangePhoneCustomer implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) ChangePhoneCustomer(ctx context.Context, userId string, data *entities.UserEntity) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.WithContext(ctx).Preload("Role").First(&user, "id = ?", userId).Where("role").Error; err != nil {
		return nil, err
	}

	user.Phone = data.Phone

	if err := u.db.WithContext(ctx).Updates(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateCustomerPoint implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) CreateCustomerPoint(ctx context.Context, userPoint *entities.UserPointEntity) error {
	err := u.db.WithContext(ctx).Create(userPoint).Error
	if err != nil {
		return err
	}

	return nil
}

// AddDetailVoucher implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) AddDetailVoucher(ctx context.Context, detailVoucher *entities.UserVoucherDetailEntity) (*entities.UserVoucherDetailEntity, error) {
	err := u.db.WithContext(ctx).Create(detailVoucher)
	if err != nil {
		return nil, err.Error
	}

	return detailVoucher, nil
}

// AddVoucher implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) AddVoucher(ctx context.Context, userVoucher *entities.UserVoucherEntity) (*entities.UserVoucherEntity, error) {
	result := u.db.WithContext(ctx).Create(userVoucher)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reload dengan preload relasi
	if err := u.db.WithContext(ctx).
		Preload("User").
		Preload("User.Role").
		Preload("Voucher").
		Preload("Voucher.User").
		Preload("Voucher.User.Role").
		Preload("Detail").
		First(userVoucher, "id = ?", userVoucher.ID).Error; err != nil {
		return nil, err
	}

	return userVoucher, nil
}

// GetCustomerPoint implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) GetCustomerPoint(ctx context.Context, customerID uuid.UUID) (*entities.UserPointEntity, error) {
	var cp entities.UserPointEntity
	if err := u.db.WithContext(ctx).First(&cp, "user_id = ?", customerID).Error; err != nil {
		return nil, err
	}

	return &cp, nil
}

// UpdateCustomerPoint implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) UpdateCustomerPoint(ctx context.Context, customerID uuid.UUID, points int) error {
	return u.db.WithContext(ctx).Model(&entities.UserPointEntity{}).
		Where("user_id = ?", customerID).
		Update("total_points", gorm.Expr("total_points - ?", points)).Error
}

// CreatePointHistory implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) CreatePointHistory(ctx context.Context, history *entities.UserPointHistoriesEntity) error {
	return u.db.WithContext(ctx).Create(history).Error
}

// helper di repo file yang sama
func dateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	// pakai UTC agar konsisten
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func (u *CustomerRepositoryImpl) FindUserPoint(
	ctx context.Context,
	customerID uuid.UUID,
) (*entities.UserPointEntity, int, *time.Time, error) {
	// ambil user_points
	var userPoint entities.UserPointEntity
	if err := u.db.WithContext(ctx).
		Preload("User").
		Preload("User.Role").
		Where("user_id = ?", customerID).
		First(&userPoint).Error; err != nil {
		return nil, 0, nil, err
	}

	// total aktif (in - out) untuk tipe yang diizinkan
	var totalPoints int
	if err := u.db.WithContext(ctx).
		Model(&entities.UserPointHistoriesEntity{}).
		Select("COALESCE(SUM(CASE WHEN direction = 'in' THEN points ELSE -points END), 0)").
		Where("user_id = ?", customerID).
		Where("point_type IN ('reward','purchase','exchange')").
		Scan(&totalPoints).Error; err != nil {
		return nil, 0, nil, err
	}
	userPoint.TotalPoints = totalPoints

	// ambil batch IN (urutkan expired_at paling awal)
	var ins []entities.UserPointHistoriesEntity
	if err := u.db.WithContext(ctx).
		Where("user_id = ? AND direction = 'in' AND point_type IN ('reward','purchase','exchange')", customerID).
		// Postgres mendukung NULLS LAST. Jika MySQL, ganti dengan "IS NULL" trick.
		Order("expired_at ASC NULLS LAST, created_at ASC").
		Find(&ins).Error; err != nil {
		return nil, 0, nil, err
	}

	// ambil batch OUT
	var outs []entities.UserPointHistoriesEntity
	if err := u.db.WithContext(ctx).
		Where("user_id = ? AND direction = 'out' AND point_type IN ('reward','purchase','exchange')", customerID).
		Order("created_at ASC").
		Find(&outs).Error; err != nil {
		return nil, 0, nil, err
	}

	// FIFO: hitung sisa poin untuk tiap batch IN tanpa mengubah nilai aslinya
	remaining := make([]int, len(ins))
	for i := range ins {
		remaining[i] = ins[i].Points
	}
	for _, o := range outs {
		toDeduct := o.Points
		for j := range remaining {
			if toDeduct <= 0 {
				break
			}
			if remaining[j] > 0 {
				if remaining[j] > toDeduct {
					remaining[j] -= toDeduct
					toDeduct = 0
				} else {
					toDeduct -= remaining[j]
					remaining[j] = 0
				}
			}
		}
	}

	// cari expired date terdekat dari batch IN yang masih punya sisa poin
	var nearest *time.Time
	for i := range ins {
		if ins[i].ExpiredAt == nil || remaining[i] <= 0 {
			continue
		}
		d := dateOnly(*ins[i].ExpiredAt)
		if nearest == nil || d.Before(*nearest) {
			tmp := d
			nearest = &tmp
		}
	}

	if nearest == nil {
		// tidak ada batch IN dengan expired_at dan sisa poin
		return &userPoint, 0, nil, nil
	}

	// total poin yang akan expired pada tanggal terdekat itu (berdasarkan sisa)
	totalExpiring := 0
	for i := range ins {
		if ins[i].ExpiredAt == nil || remaining[i] <= 0 {
			continue
		}
		if sameDay(dateOnly(*ins[i].ExpiredAt), *nearest) {
			totalExpiring += remaining[i]
		}
	}

	return &userPoint, totalExpiring, nearest, nil
}

func (u *CustomerRepositoryImpl) MarkPointsAsExpired(ctx context.Context, userID uuid.UUID, expiredDate time.Time) error {
	return u.db.WithContext(ctx).
		Model(&entities.UserPointHistoriesEntity{}).
		Where("user_id = ? AND direction = 'in' AND expired_at <= ? AND point_type IN ('reward','expired','exchange')", userID, expiredDate).
		Update("point_type", "expired").Error
}

func (u *CustomerRepositoryImpl) DecreaseTotalPoints(ctx context.Context, customerID uuid.UUID, amount int) error {
	return u.db.WithContext(ctx).
		Model(&entities.UserPointEntity{}).
		Where("user_id = ?", customerID).
		Update("total_points", gorm.Expr("total_points - ?", amount)).Error
}

// FindUserPointHistory implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindUserPointHistory(ctx context.Context, customerID uuid.UUID) ([]*entities.UserPointHistoriesEntity, error) {
	var histories []*entities.UserPointHistoriesEntity
	if err := u.db.WithContext(ctx).Preload("User").Preload("User.Role").Where("user_id = ?", customerID).Order("created_at DESC").Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}

// FindUserVoucherClaimed implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindUserVoucherClaimed(ctx context.Context, userId uuid.UUID) ([]*entities.UserVoucherEntity, error) {
	var vouchers []*entities.UserVoucherEntity
	if err := u.db.WithContext(ctx).Preload("User").Preload("User.Role").
		Preload("Voucher").
		Preload("Voucher.User").
		Preload("Voucher.User.Role").
		Preload("Detail").
		Where("user_id = ?", userId).Order("created_at DESC").Find(&vouchers).Error; err != nil {
		return nil, err
	}

	return vouchers, nil
}
