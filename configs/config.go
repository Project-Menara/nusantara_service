package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not loaded")
	}
}

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}
