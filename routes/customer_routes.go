package routes

import (
	"nusantara_service/internal/data/repositories"
	"nusantara_service/internal/domain/usecases"
	"nusantara_service/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CustomerRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	custRepo := repositories.NewCustomerRepositoryImpl(db)
	custService := usecases.NewCustomerUsecase(custRepo, rdb)
	custHandler := handlers.NewCustomerHandler(custService)

	e.POST("/check-phone", custHandler.CheckPhoneCustomer)
	e.POST("/register", custHandler.RegisterCustomer)
	e.POST("/resend-code-verify", custHandler.ResendCodeOTPVerify)
}
