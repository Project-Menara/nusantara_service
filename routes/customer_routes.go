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

func CustomerRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client, cloudinarySvc *cloudinary.CloudinaryService) {
	custRepo := repositories.NewCustomerRepositoryImpl(db)
	voucherRepo := repositories.NewVoucherRepositoryImpl(db)
	custService := usecases.NewCustomerUsecase(custRepo, rdb, cloudinarySvc, db, voucherRepo)
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
	e.PUT("/update-profile", custHandler.UpdateProfile, middlewares.JWTMiddleware(rdb))
	e.POST("/verify-pin", custHandler.VerifyPINCustomer, middlewares.JWTMiddleware(rdb))
	e.POST("/new-pin-customer", custHandler.NewPinCustomer, middlewares.JWTMiddleware(rdb))
	e.PUT("/confirm-new-pin-customer", custHandler.ConfirmationPINCustomer, middlewares.JWTMiddleware(rdb))
	e.POST("/new-phone-customer", custHandler.NewPhoneCustomer, middlewares.JWTMiddleware(rdb))
	e.PUT("/verify-otp-customer", custHandler.VerifyCodeOTPCustomerUpdate, middlewares.JWTMiddleware(rdb))
	e.POST("/forgot-pin", custHandler.ForgotPINCustomer)
	e.GET("/validate-forgot-pin", custHandler.ValidateTokenForgotPINCustomer)
	e.POST("/new-pin-forgot", custHandler.NewPINCustomerForgotPIN)
	e.POST("/confirm-pin-forgot", custHandler.ConfirmationPINCustomerForgotPIN)

	e.POST("/claim-voucher/:voucherId", custHandler.ClaimVoucherCustomer, middlewares.JWTMiddleware(rdb))
	e.GET("/point", custHandler.GetCustomerPoint, middlewares.JWTMiddleware(rdb))
	e.GET("/point/history", custHandler.GetCustomerPointHistory, middlewares.JWTMiddleware(rdb))
	e.GET("/vouchers/claimed", custHandler.GetCustomerVouchersClaimed, middlewares.JWTMiddleware(rdb))

	e.GET("/shop-detail/:shop_id", custHandler.GetShopByID)
	e.GET("/my-cart", custHandler.GetMyCart, middlewares.JWTMiddleware(rdb))
	e.POST("/add-cart-item", custHandler.AddProductToMyCart, middlewares.JWTMiddleware(rdb))
	e.DELETE("/delete-cart-item/:product_id", custHandler.DeleteMyCartItem, middlewares.JWTMiddleware(rdb))
}
