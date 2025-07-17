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

func UserRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	userRepo := repositories.NewUserRepositoryImpl(db)
	authService := usecases.NewAuthUsecase(userRepo, rdb)
	authHandler := handlers.NewAuthHandler(authService)

	e.POST("/register", authHandler.RegisterUser)
	e.POST("/login", authHandler.LoginUser)
	e.POST("/logout", authHandler.LogoutUser, middlewares.JWTMiddleware(rdb))
	e.GET("/me", authHandler.GetProfile, middlewares.JWTMiddleware(rdb))
}
