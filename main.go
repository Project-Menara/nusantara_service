package main

import (
	"log"
	"nusantara_service/configs"
	"nusantara_service/internal/middlewares"
	"nusantara_service/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	configs.LoadEnv()

	db := configs.InitDB()
	rdb := configs.InitRedis()

	configs.RunMigrations(db)

	e := echo.New()

	e.Use(middlewares.LoggerMiddleware)

	routes.Routes(e, db, rdb)

	log.Fatal(e.Start(":8080"))
}
