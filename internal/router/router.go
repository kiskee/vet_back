package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/danielm/app_sara_backend/internal/auth"
	"github.com/danielm/app_sara_backend/internal/middleware"
	userHandlerPkg "github.com/danielm/app_sara_backend/internal/user"

	"github.com/danielm/app_sara_backend/internal/domain"
)

type Dependencies struct {
	AuthHandler      *auth.Handler
	UserHandler      *userHandlerPkg.Handler
	JWTSecret        string
}

func Setup(app *fiber.App, deps *Dependencies) {
	app.Use(logger.New())

	api := app.Group("/api/v1")
	authMiddleware := middleware.AuthRequired(deps.JWTSecret)
	adminMiddleware := middleware.RoleRequired(domain.RoleAdmin)

	auth.RegisterRoutes(api, deps.AuthHandler)
	userHandlerPkg.RegisterRoutes(api, deps.UserHandler, authMiddleware, adminMiddleware)
}
