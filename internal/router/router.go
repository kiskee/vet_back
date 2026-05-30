package router

import (
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/danielm/app_sara_backend/internal/auth"
	"github.com/danielm/app_sara_backend/internal/middleware"
	servicesPkg "github.com/danielm/app_sara_backend/internal/services"
	userHandlerPkg "github.com/danielm/app_sara_backend/internal/user"
	wsapp "github.com/danielm/app_sara_backend/internal/websocket"

	"github.com/danielm/app_sara_backend/internal/domain"
)

type Dependencies struct {
	AuthHandler    *auth.Handler
	UserHandler    *userHandlerPkg.Handler
	ServicesHandler *servicesPkg.Handler
	WSHandler      *wsapp.Handler
	JWTSecret      string
}

func Setup(app *fiber.App, deps *Dependencies) {
	app.Use(logger.New())

	api := app.Group("/api/v1")
	authMiddleware := middleware.AuthRequired(deps.JWTSecret)
	adminMiddleware := middleware.RoleRequired(domain.RoleAdmin)

	auth.RegisterRoutes(api, deps.AuthHandler,
		middleware.RateLimit(10, time.Minute),
	)
	userHandlerPkg.RegisterRoutes(api, deps.UserHandler, authMiddleware, adminMiddleware,
		middleware.RateLimit(30, time.Minute),
	)
	servicesPkg.RegisterRoutes(api, deps.ServicesHandler, authMiddleware, adminMiddleware,
		middleware.RateLimit(30, time.Minute),
	)

	wsAuth := wsapp.AuthMiddlewareForWS(deps.JWTSecret)
	app.Get("/ws/client", wsAuth, websocket.New(deps.WSHandler.HandleClient))
	app.Get("/ws/vet", wsAuth, websocket.New(deps.WSHandler.HandleVet))
}
