package repositories

import (
	"context"
	"errors"
	"nusantara_service/internal/data/model"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	cartresponse "nusantara_service/internal/dto/responses/cart_response"
	favoriteresponse "nusantara_service/internal/dto/responses/favorite_response"
	shopresponse "nusantara_service/internal/dto/responses/shop_response"
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

	now := time.Now().UTC()

	// 1) AUTO-EXPIRE: ubah point_type ke 'expired' untuk IN yang sudah lewat expired_at.
	//    Gunakan waktu UTC supaya konsisten.
	if err := u.db.WithContext(ctx).
		Model(&entities.UserPointHistoriesEntity{}).
		Where(`
			user_id = ? 
			AND direction = 'in' 
			AND point_type IN ('reward','purchase','exchange') 
			AND expired_at IS NOT NULL 
			AND expired_at < ?
		`, customerID, now).
		Updates(map[string]any{
			"point_type": "expired",
			"updated_at": now,
		}).Error; err != nil {
		return nil, 0, nil, err
	}

	// 2) Ambil parent user_points (beserta user & role)
	var userPoint entities.UserPointEntity
	if err := u.db.WithContext(ctx).
		Preload("User").
		Preload("User.Role").
		Where("user_id = ?", customerID).
		First(&userPoint).Error; err != nil {
		return nil, 0, nil, err
	}

	// 3) Ambil seluruh batch IN (termasuk yang sudah ditandai 'expired') dan semua OUT aktif.
	//    Kita butuh semua IN agar OUT yang terjadi SEBELUM expired tetap bisa mengonsumsi batch tersebut.
	var ins []entities.UserPointHistoriesEntity
	if err := u.db.WithContext(ctx).
		Where(`
			user_id = ? 
			AND direction = 'in' 
			AND point_type IN ('reward','purchase','exchange','expired')
		`, customerID).
		Order("created_at ASC").
		Find(&ins).Error; err != nil {
		return nil, 0, nil, err
	}

	var outs []entities.UserPointHistoriesEntity
	if err := u.db.WithContext(ctx).
		Where(`
			user_id = ? 
			AND direction = 'out' 
			AND point_type IN ('reward','purchase','exchange')
		`, customerID).
		Order("created_at ASC").
		Find(&outs).Error; err != nil {
		return nil, 0, nil, err
	}

	// 4) Rekonstruksi FIFO berbasis waktu.
	type batch struct {
		remaining int
		expiredAt *time.Time
		createdAt time.Time
	}

	batches := make([]batch, len(ins))
	for i := range ins {
		batches[i] = batch{
			remaining: ins[i].Points,
			expiredAt: ins[i].ExpiredAt,
			createdAt: ins[i].CreatedAt,
		}
	}

	// helper: expire semua batch yang expired_at <= t (inklusif)
	expireUpTo := func(t time.Time) {
		for i := range batches {
			if batches[i].remaining <= 0 || batches[i].expiredAt == nil {
				continue
			}
			if !batches[i].expiredAt.After(t) {
				// expired_at <= t  → hanguskan sisa batch
				batches[i].remaining = 0
			}
		}
	}

	// Proses semua OUT secara kronologis.
	for _, o := range outs {
		// 4a) Sebelum OUT terjadi, hanguskan batch yang sudah expired sampai waktu OUT
		expireUpTo(o.CreatedAt)

		toDeduct := o.Points

		// 4b) Deduct hanya dari batch yang SUDAH ADA (created_at <= o.CreatedAt)
		for j := range batches {
			if toDeduct <= 0 {
				break
			}
			// batch dari masa depan tidak boleh dipakai
			if batches[j].createdAt.After(o.CreatedAt) {
				continue
			}
			if batches[j].remaining > 0 {
				if batches[j].remaining > toDeduct {
					batches[j].remaining -= toDeduct
					toDeduct = 0
				} else {
					toDeduct -= batches[j].remaining
					batches[j].remaining = 0
				}
			}
		}
		// Catatan: sisa toDeduct (jika > 0) diabaikan → tidak meminjam dari masa depan.
	}

	// 4c) Setelah semua OUT, hanguskan batch yang expired sampai "sekarang"
	expireUpTo(now)

	// 5) Hitung total aktif saat ini + nearest expiry
	totalPoints := 0
	var nearest *time.Time
	for i := range batches {
		if batches[i].remaining <= 0 {
			continue
		}
		totalPoints += batches[i].remaining

		if batches[i].expiredAt != nil && batches[i].expiredAt.After(now) {
			d := dateOnly(*batches[i].expiredAt)
			if nearest == nil || d.Before(*nearest) {
				tmp := d
				nearest = &tmp
			}
		}
	}

	// 6) Berapa yang akan expired pada tanggal terdekat itu?
	totalExpiring := 0
	if nearest != nil {
		for i := range batches {
			if batches[i].remaining <= 0 || batches[i].expiredAt == nil {
				continue
			}
			if sameDay(dateOnly(*batches[i].expiredAt), *nearest) {
				totalExpiring += batches[i].remaining
			}
		}
	}

	// 7) Sinkronkan kolom total_points di tabel user_points
	if err := u.db.WithContext(ctx).
		Model(&entities.UserPointEntity{}).
		Where("user_id = ?", customerID).
		Update("total_points", totalPoints).Error; err != nil {
		return nil, 0, nil, err
	}
	userPoint.TotalPoints = totalPoints

	// 8) Return
	if nearest == nil {
		return &userPoint, 0, nil, nil
	}
	return &userPoint, totalExpiring, nearest, nil
}

func (u *CustomerRepositoryImpl) MarkPointsAsExpired(ctx context.Context, userID uuid.UUID, expiredDate time.Time) error {
	return u.db.WithContext(ctx).
		Model(&entities.UserPointHistoriesEntity{}).
		Where("user_id = ? AND direction = 'in' AND expired_at <= ? AND point_type IN ('reward','exchange')", userID, expiredDate).
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

// FindShopById implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) FindShopById(ctx context.Context, offset int, limit int, search string, typeID uuid.UUID, shopID uuid.UUID) (*shopresponse.ShopCustomerResponse, int, error) {
	var shop entities.ShopEntity
	var totalProducts int64
	countQuery := u.db.WithContext(ctx).Model(&entities.ShopProductEntity{}).
		Where("shop_id = ?", shopID)

	if search != "" || typeID != uuid.Nil {
		countQuery = countQuery.
			Joins("JOIN products ON shop_products.product_id = products.id")
	}

	if search != "" {
		countQuery = countQuery.Where("products.name ILIKE ?", "%"+search+"%")
	}
	if typeID != uuid.Nil {
		countQuery = countQuery.Where("products.type_product_id = ?", typeID)
	}

	if err := countQuery.Count(&totalProducts).Error; err != nil {
		return nil, 0, err
	}

	productScope := func(db *gorm.DB) *gorm.DB {
		if search != "" || typeID != uuid.Nil {
			db = db.
				Joins("JOIN products ON products.id = shop_products.product_id")
		}

		if search != "" {
			db = db.Where("products.name ILIKE ?", "%"+search+"%")
		}
		if typeID != uuid.Nil {
			db = db.Where("products.type_product_id = ?", typeID)
		}

		if limit > 0 {
			db = db.Limit(limit).Offset(offset)
		}
		return db
	}

	if err := u.db.WithContext(ctx).
		Preload("ShopImages.Image").
		Preload("ShopProducts", productScope).
		Preload("ShopProducts.Product").
		Preload("ShopProducts.Product.Image").
		Preload("ShopProducts.Product.ProductImages.Image").
		Preload("ShopProducts.Product.TypeProduct").
		Preload("ShopProducts.Product.User").
		First(&shop, "id = ?", shopID).Error; err != nil {

		return nil, 0, err
	}

	responses := shopresponse.ToShopCustomerResponse(shop)

	return &responses, int(totalProducts), nil
}

// GetMyCart implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) GetMyCart(ctx context.Context, customerID uuid.UUID) (*cartresponse.CartResponse, error) {
	var cart entities.CartEntity

	err := u.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", customerID, 0).
		Preload("CartItems.Product").
		Preload("CartItems.Product.Image").
		Preload("CartItems.Product.ProductImages.Image").
		Preload("CartItems.Product.TypeProduct").
		Preload("CartItems.Product.User").
		Preload("User").
		Preload("Shop").
		First(&cart).Error

	if err != nil {
		return nil, err
	}

	cartResponse := cartresponse.ToCartResponse(cart)

	return &cartResponse, nil
}

// CreateMyCart implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) CreateMyCart(ctx context.Context, cart *entities.CartEntity) (*entities.CartEntity, error) {
	if err := u.db.WithContext(ctx).Create(cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

// AddProductToCart implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) AddProductToCart(ctx context.Context, cartID uuid.UUID, productID uuid.UUID) error {
	var cartItem entities.CartItemEntity

	res := u.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		First(&cartItem)

	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	}

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		cartItem = entities.CartItemEntity{
			ID:        uuid.New(),
			CartID:    cartID,
			ProductID: productID,
			Selected:  false,
		}
		return u.db.WithContext(ctx).Create(&cartItem).Error
	}

	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	return tx.Commit().Error
}

// UpdateSelectedCartItem implements [repositories.CustomerRepository].
func (u *CustomerRepositoryImpl) UpdateSelectedCartItem(ctx context.Context, cartID uuid.UUID, productID uuid.UUID, selected bool) error {
	query := u.db.WithContext(ctx).Model(&model.CartItem{}).Where("cart_id = ?", cartID)

	// Jika ingin update spesifik produk
	if productID != uuid.Nil {
		query = query.Where("product_id = ?", productID)
	}

	result := query.Update("selected", selected)

	if result.Error != nil {
		return result.Error
	}

	// Jika tidak ada baris yang ter-update, artinya ID tidak ditemukan
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteCartItem implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) DeleteCartItem(ctx context.Context, cartID uuid.UUID, productID uuid.UUID) error {
	res := u.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Delete(&entities.CartItemEntity{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetMyFavorite implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) GetMyFavorite(ctx context.Context, customerID uuid.UUID) (*favoriteresponse.FavoriteResponse, error) {
	var favorite entities.FavoriteEntity

	err := u.db.WithContext(ctx).
		Where("user_id = ?", customerID).
		Preload("FavoriteItems.Product").
		Preload("FavoriteItems.Product.Image").
		Preload("FavoriteItems.Product.ProductImages.Image").
		Preload("FavoriteItems.Product.TypeProduct").
		Preload("FavoriteItems.Product.User").
		Preload("User").
		First(&favorite).Error

	if err != nil {
		return nil, err
	}

	favoriteResponse := favoriteresponse.ToFavoriteResponse(favorite)
	return &favoriteResponse, nil
}

// CreateMyFavorite implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) CreateMyFavorite(ctx context.Context, favorite *entities.FavoriteEntity) (*entities.FavoriteEntity, error) {
	if err := u.db.WithContext(ctx).Create(favorite).Error; err != nil {
		return nil, err
	}
	return favorite, nil
}

// AddProductToFavorite implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) AddProductToFavorite(ctx context.Context, favoriteID uuid.UUID, productID uuid.UUID) error {
	var favoriteItem entities.FavoriteItemEntity

	res := u.db.WithContext(ctx).
		Where("favorite_id = ? AND product_id = ?", favoriteID, productID).
		First(&favoriteItem)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	}

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		favoriteItem = entities.FavoriteItemEntity{
			ID:         uuid.New(),
			FavoriteID: favoriteID,
			ProductID:  productID,
			Selected:   true,
		}
		return u.db.WithContext(ctx).Create(&favoriteItem).Error
	}

	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	return tx.Commit().Error
}

// DeleteFavoriteItem implements repositories.CustomerRepository.
func (u *CustomerRepositoryImpl) DeleteFavoriteItem(ctx context.Context, favoriteID uuid.UUID, productID uuid.UUID) error {
	res := u.db.WithContext(ctx).
		Where("favorite_id = ? AND product_id = ?", favoriteID, productID).
		Delete(&entities.FavoriteItemEntity{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
