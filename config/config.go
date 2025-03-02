package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	JWTSecret    string
	UploadsDir   string
	DownloadsDir string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default variables")
	}

	return &Config{
		AppPort:      getEnv("APP_PORT", "3000"),
		JWTSecret:    getEnv("JWT_SECRET", "jwt_secret"),
		UploadsDir:   getEnv("UPLOADS_DIR", "./media/uploads/"),
		DownloadsDir: getEnv("DOWNLOADS_DIR", "./media/downloads/"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
