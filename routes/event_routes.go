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

func EventRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client, cld *cloudinary.CloudinaryService) {
	eventRepo := repositories.NewEventRepositoryImpl(db)
	productRepo := repositories.NewProductRepositoryImpl(db)
	eventService := usecases.NewEventUsecase(eventRepo, productRepo, rdb, *cld)
	eventHandler := handlers.NewEventHandler(eventService)

	e.POST("/create", eventHandler.CreateEvent, middlewares.JWTMiddleware(rdb))
	e.GET("", eventHandler.GetAllEvents, middlewares.JWTMiddleware(rdb))
	e.GET("/:id", eventHandler.GetEventById, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit", eventHandler.UpdateEvent, middlewares.JWTMiddleware(rdb))
	e.DELETE("/:id/delete", eventHandler.DeleteEvent, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit-status", eventHandler.UpdateStatusEvent, middlewares.JWTMiddleware(rdb))
}
