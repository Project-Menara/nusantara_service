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

func BannerRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client, cloudinary *cloudinary.CloudinaryService) {
	bannerRepo := repositories.NewBannerRepositoryImpl(db)
	bannerService := usecases.NewBannerUsecase(bannerRepo, rdb, cloudinary)
	bannerHandler := handlers.NewBannerHandler(bannerService)

	e.POST("/create", bannerHandler.CreateBanner, middlewares.JWTMiddleware(rdb))
	e.GET("", bannerHandler.GetAllBanner, middlewares.JWTMiddleware(rdb))
	e.GET("/:id", bannerHandler.GetByIdBanner, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit", bannerHandler.UpdateBanner, middlewares.JWTMiddleware(rdb))
	e.DELETE("/:id/delete", bannerHandler.DeleteBanner, middlewares.JWTMiddleware(rdb))
}
