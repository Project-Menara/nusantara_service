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

func ShopRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client, cld *cloudinary.CloudinaryService) {
	prodRepo := repositories.NewProductRepositoryImpl(db)
	shopRepo := repositories.NewShopRepositoryImpl(db)
	shopService := usecases.NewShopUsecase(shopRepo, rdb, cld, prodRepo)
	shopHandler := handlers.NewShopHandler(shopService)

	e.POST("/create", shopHandler.CreateShop, middlewares.JWTMiddleware(rdb))
	e.GET("", shopHandler.GetAllShop, middlewares.JWTMiddleware(rdb))
	e.GET("/:id", shopHandler.GetByIdShop, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit", shopHandler.UpdateShop, middlewares.JWTMiddleware(rdb))
	e.DELETE("/:id/delete", shopHandler.DeleteShop, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit-status", shopHandler.UpdateStatusShop, middlewares.JWTMiddleware(rdb))
}
