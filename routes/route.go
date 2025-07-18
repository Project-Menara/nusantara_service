package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Routes(e *echo.Echo, db *gorm.DB, rdb *redis.Client) {
	v1 := e.Group("/api/v1")
	RoleRoutes(v1.Group("/role"), db, rdb)
	UserRoutes(v1.Group("/user"), db, rdb)
}
