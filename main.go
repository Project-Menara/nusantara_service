package main

import (
	"log"
	"nusantara_service/configs"
	"nusantara_service/internal/middlewares"
	"nusantara_service/routes"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	configs.LoadEnv()

	db := configs.InitDB()
	rdb := configs.InitRedis()

	configs.RunMigrations(db)

	e := echo.New()

	e.Use(middlewares.LoggerMiddleware)

	for _, r := range e.Routes() {
		log.Printf("ROUTE %s %s", r.Method, r.Path)
	}

	routes.Routes(e, db, rdb)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))

}
