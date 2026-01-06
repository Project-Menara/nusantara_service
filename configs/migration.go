package configs

import (
	"nusantara_service/internal/data/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.Banner{},
		&model.TypeProduct{},
		&model.UserPoint{},
		&model.UserPointHistories{},
		&model.Voucher{},
		&model.UserVoucher{},
		&model.UserVoucherDetail{},
		&model.Image{},
		&model.Product{},
		&model.ProductImage{},
		&model.Shop{},
		&model.ShopCashier{},
		&model.ShopProducts{},
		&model.ShopImage{},
		&model.CustomerAddress{},
		&model.Event{},
		&model.EventProduct{},
		&model.EventBundleBuy{},
		&model.EventBundleReward{},
		&model.Cart{},
		&model.CartItem{},
		&model.Favorite{},
		&model.FavoriteItem{},
		&model.Order{},
		&model.OrderItem{},
		&model.OrderEvent{},
		&model.OrderReward{},
		&model.OrderVoucher{},
	)
}
