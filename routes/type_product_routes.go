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

func TypeProductRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client, cloudinary *cloudinary.CloudinaryService) {
	typeProductRepo := repositories.NewTypeProductRepositoryImpl(db)
	typeProductService := usecases.NewTypeProductUsecase(typeProductRepo, rdb, cloudinary)
	typeProductHandler := handlers.NewTypeProductHandler(typeProductService)

	e.POST("/create", typeProductHandler.CreateTypeProduct, middlewares.JWTMiddleware(rdb))
	e.GET("", typeProductHandler.GetAllTypeProduct, middlewares.JWTMiddleware(rdb))
	e.GET("/:id", typeProductHandler.GetByIdTypeProduct, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit", typeProductHandler.UpdateTypeProduct, middlewares.JWTMiddleware(rdb))
	e.DELETE("/:id/delete", typeProductHandler.DeleteTypeProduct, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit-status", typeProductHandler.UpdateStatusTypeProduct, middlewares.JWTMiddleware(rdb))

	e.GET("/customer", typeProductHandler.GetAllTypeProductCustomer)
	e.GET("/:id/customer", typeProductHandler.GetByIdTypeProductCustomer)
}
