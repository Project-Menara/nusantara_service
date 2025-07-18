package routes

import (
	"nusantara_service/internal/data/repositories"
	"nusantara_service/internal/domain/usecases"
	"nusantara_service/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RoleRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	roleRepo := repositories.NewRoleRepositoryImpl(db)
	roleService := usecases.NewRoleUsecase(roleRepo, rdb)
	roleHandler := handlers.NewRoleHandler(roleService)

	e.POST("/create", roleHandler.CreateRole)
}
