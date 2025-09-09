package routes

import (
	"nusantara_service/internal/data/dataSources/cloudinary"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Routes(e *echo.Echo, db *gorm.DB, rdb *redis.Client, cloudinarySvc *cloudinary.CloudinaryService) {
	v1 := e.Group("/api/v1")
	RoleRoutes(v1.Group("/role"), db, rdb)
	UserRoutes(v1.Group("/user"), db, rdb)
	CustomerRoutes(v1.Group("/customer"), db, rdb, cloudinarySvc)
	BannerRoutes(v1.Group("/banner"), db, rdb, cloudinarySvc)
	TypeProductRoutes(v1.Group("/type-product"), db, rdb, cloudinarySvc)
	VoucherRoutes(v1.Group("/voucher"), db, rdb, cloudinarySvc)
	ProductRoutes(v1.Group("/product"), db, rdb, cloudinarySvc)
	CashierRoutes(v1.Group("/cashier"), db, rdb, cloudinarySvc)
	ShopRoutes(v1.Group("/shop"), db, rdb, cloudinarySvc)
}
