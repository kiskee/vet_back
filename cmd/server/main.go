package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"github.com/danielm/app_sara_backend/internal/auth"
	"github.com/danielm/app_sara_backend/internal/config"
	"github.com/danielm/app_sara_backend/internal/database"
	"github.com/danielm/app_sara_backend/internal/router"
	servicesPkg "github.com/danielm/app_sara_backend/internal/services"
	userHandlerPkg "github.com/danielm/app_sara_backend/internal/user"
	vetsPkg "github.com/danielm/app_sara_backend/internal/vets"
	"github.com/danielm/app_sara_backend/internal/websocket"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("connected to database")

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	defer sqlDB.Close()

	rdb, err := database.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer rdb.Close()
	log.Println("connected to redis")

	userRepo := userHandlerPkg.NewRepository(db)
	userService := userHandlerPkg.NewService(userRepo)
	userHandler := userHandlerPkg.NewHandler(userService)

	vetsRepo := vetsPkg.NewRepository(db, rdb)
	vetsService := vetsPkg.NewService(vetsRepo, rdb)
	vetsHandler := vetsPkg.NewHandler(vetsService)

	authService := auth.NewService(userRepo, vetsRepo, cfg.JWTSecret, cfg.JWTRefreshSecret, cfg.AdminSecret)
	authHandler := auth.NewHandler(authService)

	servicesRepo := servicesPkg.NewRepository(db)
	servicesService := servicesPkg.NewService(servicesRepo)
	servicesHandler := servicesPkg.NewHandler(servicesService)

	wsHub := websocket.NewHub()
	wsHandler := websocket.NewHandler(wsHub)

	app := fiber.New(fiber.Config{
		AppName: "app_sara_backend",
	})

	deps := &router.Dependencies{
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		ServicesHandler: servicesHandler,
		VetsHandler:     vetsHandler,
		WSHandler:       wsHandler,
		JWTSecret:       cfg.JWTSecret,
	}

	router.Setup(app, deps)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
	log.Printf("server listening on port %s", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
}
