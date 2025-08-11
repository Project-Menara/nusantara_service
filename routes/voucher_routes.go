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

func VoucherRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client, cloudinary *cloudinary.CloudinaryService) {
	voucherRepo := repositories.NewVoucherRepositoryImpl(db)
	voucherService := usecases.NewVoucherUsecase(voucherRepo, rdb, cloudinary)
	voucherHandler := handlers.NewVoucherHandler(voucherService)

	e.POST("/create", voucherHandler.CreateVoucher, middlewares.JWTMiddleware(rdb))
	e.GET("", voucherHandler.GetAllVoucher, middlewares.JWTMiddleware(rdb))
	e.GET("/:id", voucherHandler.GetByIdVoucher, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit", voucherHandler.UpdateVoucher, middlewares.JWTMiddleware(rdb))
	e.DELETE("/:id/delete", voucherHandler.DeleteVoucher, middlewares.JWTMiddleware(rdb))
	e.PUT("/:id/edit-status", voucherHandler.UpdateStatusVoucher, middlewares.JWTMiddleware(rdb))

	e.GET("/customer", voucherHandler.GetAllVoucherCustomer, middlewares.JWTMiddleware(rdb))
	e.GET("/:id/customer", voucherHandler.GetByIdVoucherCustomer, middlewares.JWTMiddleware(rdb))
}
