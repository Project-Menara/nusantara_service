package configs

import (
	"nusantara_service/internal/data/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
	)
}
