package routes

import (
	"nusantara_service/internal/data/repositories"
	"nusantara_service/internal/domain/usecases"
	"nusantara_service/internal/handlers"
	"nusantara_service/internal/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CustomerAddressRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	custAddrRepo := repositories.NewCustomerAddresRepositoryImpl(db)
	custAddrService := usecases.NewCustomerAddressUsecase(custAddrRepo, rdb)
	custAddrHandler := handlers.NewCustomerAddressHandler(custAddrService)

	// Routes
	e.POST("/addresses/create", custAddrHandler.CreateAddress, middlewares.JWTMiddleware(rdb))             // Create address
	e.GET("/addresses", custAddrHandler.GetAllAddress, middlewares.JWTMiddleware(rdb))                     // Get all addresses
	e.GET("/addresses/:id", custAddrHandler.GetByIdAddress, middlewares.JWTMiddleware(rdb))                // Get address by ID
	e.PUT("/addresses/:id/edit", custAddrHandler.UpdateAddress, middlewares.JWTMiddleware(rdb))            // Update address
	e.DELETE("/addresses/:id/delete", custAddrHandler.DeleteAddress, middlewares.JWTMiddleware(rdb))       // Delete address
	e.GET("/addresses/default", custAddrHandler.GetDefaultAddress, middlewares.JWTMiddleware(rdb))         // Get default address
	e.PUT("/addresses/:id/set-default", custAddrHandler.SetDefaultAddress, middlewares.JWTMiddleware(rdb)) // Set default address
	e.GET("/addresses/nearby-shops", custAddrHandler.GetNearbyShops, middlewares.JWTMiddleware(rdb))       // Get nearby shops
	e.GET("/addresses/public-nearby-shops", custAddrHandler.GetNearbyShops)
}
