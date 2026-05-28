package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	RedisURL          string
	JWTSecret         string
	JWTRefreshSecret  string
	Port              string
	AdminSecret       string
}

func Load() *Config {
	godotenv.Load()

	cfg := &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/app_sara?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		Port:             getEnv("PORT", "3000"),
		AdminSecret:      os.Getenv("ADMIN_SECRET"),
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	if cfg.JWTRefreshSecret == "" {
		log.Fatal("JWT_REFRESH_SECRET is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
