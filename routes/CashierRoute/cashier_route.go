package cashierroute

import (
	cashierrepoimpl "nusantara_service/internal/data/repositories/CashierRepoImpl"
	cashierusecase "nusantara_service/internal/domain/usecases/CashierUsecase"
	cashierhandler "nusantara_service/internal/handlers/CashierHandler"
	"nusantara_service/internal/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CashierRoute(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	cashierRepo := cashierrepoimpl.NewCashierRepositoryImpl(db)
	cashierService := cashierusecase.NewCashierUsecase(cashierRepo, rdb)
	cashierHandler := cashierhandler.NewCashierHandler(cashierService)

	e.GET("/shop-names", cashierHandler.GetAllShopName, middlewares.JWTMiddleware(rdb))
	e.GET("/shop-cashier/:shop_id", cashierHandler.GetDetailShopCashier, middlewares.JWTMiddleware(rdb))
	e.GET("/cashier-shop-product/:shop_id", cashierHandler.GetShopProductCashier, middlewares.JWTMiddleware(rdb))
}
