package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	JWTSecret         string
	JWTRefreshSecret  string
	Port              string
	AdminSecret       string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/app_sara?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "secret"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "refresh-secret"),
		Port:             getEnv("PORT", "3000"),
		AdminSecret:      getEnv("ADMIN_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
