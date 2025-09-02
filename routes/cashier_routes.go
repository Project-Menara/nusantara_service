package routes

import (
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/repositories"
	"nusantara_service/internal/domain/usecases"
	"nusantara_service/internal/handlers"
	"nusantara_service/internal/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CashierRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client, cloudinary *cloudinary.CloudinaryService) {
	cashierRepo := repositories.NewCashierRepositoryImpl(db)
	roleRepo := repositories.NewRoleRepositoryImpl(db)
	cashierService := usecases.NewCashierUsecase(cashierRepo, roleRepo, rdb, cloudinary)
	cashierHandler := handlers.NewCashierHandler(cashierService)

	e.POST("/create", cashierHandler.CreateCashier, middlewares.JWTMiddleware(rdb))
	e.GET("", cashierHandler.GetCashierAll, middlewares.JWTMiddleware(rdb))
	e.GET("/:id", cashierHandler.GetCashierById, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit", cashierHandler.UpdateCashier, middlewares.JWTMiddleware(rdb))
	e.DELETE("/:id/delete", cashierHandler.DeleteCashier, middlewares.JWTMiddleware(rdb))
}
