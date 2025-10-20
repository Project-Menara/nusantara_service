package main

import (
	"log"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/middlewares"
	"nusantara_service/internal/workers/producer"
	"nusantara_service/routes"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// @title API DOCS NUSANTARA OLEH OLEH
	// @version 1.0
	// @description This is a Nusantara Oleh Oleh Apps
	// @termOfService http://swagger.io/terms/

	// @contact.name API Support
	// @contact.url http://www.swagger.io/support
	// @contact.email support@swagger.io

	// @license.name Apache 2.0
	// @license.url http://www.apache.org/license/LICENSE-2.0.html

	// @host https://nusantaraservice-production.up.railway.app
	// @BasePath /api/v1
	configs.LoadEnv()

	db := configs.InitDB()
	rdb := configs.InitRedis()

	configs.RunMigrations(db)

	e := echo.New()

	e.Use(middlewares.LoggerMiddleware)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	for _, r := range e.Routes() {
		log.Printf("ROUTE %s %s", r.Method, r.Path)
	}

	cloudinarySvc, err := cloudinary.NewCloudinaryService()
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary service: %v", err)
	}

	configs.InitRabbitMQ()
	defer configs.CloseConnections()

	go producer.StartWorkers()

	routes.Routes(e, db, rdb, &cloudinarySvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))
}
