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
	)
}
