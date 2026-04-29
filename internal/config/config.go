package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	AppEnv    string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}
	return &Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		AppEnv:    getEnv("APP_ENV", "production"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "ecom_user"),
		DBPass:    getEnv("DB_PASSWORD", "ecom_pass123"),
		DBName:    getEnv("DB_NAME", "ecom_db"),
		JWTSecret: getEnv("JWT_SECRET", "changeme"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
