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

func ProductRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client, cloudinary *cloudinary.CloudinaryService) {
	productRepo := repositories.NewProductRepositoryImpl(db)
	typeProductRepo := repositories.NewTypeProductRepositoryImpl(db)
	productService := usecases.NewProductUsecase(productRepo, typeProductRepo, rdb, cloudinary)
	productHandler := handlers.NewProductHandler(productService)

	e.POST("/create", productHandler.CreateProduct, middlewares.JWTMiddleware(rdb))
	e.GET("", productHandler.GetAllProduct, middlewares.JWTMiddleware(rdb))
	e.GET("/:id", productHandler.GetByIDProduct, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit", productHandler.UpdateProduct, middlewares.JWTMiddleware(rdb))
	e.DELETE("/:id/delete", productHandler.DeleteProduct, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit-status", productHandler.UpdateStatusProduct, middlewares.JWTMiddleware(rdb))
}
