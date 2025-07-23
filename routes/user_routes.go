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
	authService := usecases.NewUserUsecase(userRepo, rdb)
	authHandler := handlers.NewAuthHandler(authService)

	e.POST("/admin/register", authHandler.RegisterAdmin)
	e.POST("/admin/login", authHandler.LoginAdmin)
	e.GET("/admin/me", authHandler.GetProfileAdmin, middlewares.JWTMiddleware(rdb))
	e.POST("/admin/logout", authHandler.LogoutAdmin, middlewares.JWTMiddleware(rdb))
}
