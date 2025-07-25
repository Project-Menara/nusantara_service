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

func CustomerRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	custRepo := repositories.NewCustomerRepositoryImpl(db)
	custService := usecases.NewCustomerUsecase(custRepo, rdb)
	custHandler := handlers.NewCustomerHandler(custService)

	e.POST("/check-phone", custHandler.CheckPhoneCustomer)
	e.POST("/register", custHandler.RegisterCustomer)
	e.POST("/resend-code-verify", custHandler.ResendCodeOTPVerify)
	e.POST("/code-verify", custHandler.VerifyCodeOTP)
	e.POST("/new-pin", custHandler.NewPin)
	e.POST("/confirm-pin", custHandler.ConfirmationPin)
	e.GET("/me", custHandler.GetProfileCustomer, middlewares.JWTMiddleware(rdb))
	e.POST("/logout", custHandler.LogoutCustomer, middlewares.JWTMiddleware(rdb))
	e.POST("/login", custHandler.LoginCustomer)
}
