package orderroute

import (
	"nusantara_service/internal/data/repositories"
	orderrepoimpl "nusantara_service/internal/data/repositories/order_repo_impl"
	orderusecase "nusantara_service/internal/domain/usecases/order_usecase"
	orderhandler "nusantara_service/internal/handlers/order_handler"
	"nusantara_service/internal/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func OrderRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	orderRepo := orderrepoimpl.NewOrderRepositoryImpl(db)
	userRepo := repositories.NewUserRepositoryImpl(db)
	productRepo := repositories.NewProductRepositoryImpl(db)
	orderService := orderusecase.NewOrderUsecase(orderRepo, userRepo, productRepo, rdb)
	orderHandler := orderhandler.NewOrderHandler(orderService)

	e.POST("/create", orderHandler.CreateOrder, middlewares.JWTMiddleware(rdb))
}
