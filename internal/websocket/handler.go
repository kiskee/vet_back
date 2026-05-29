package websocket

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"github.com/danielm/app_sara_backend/internal/domain"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) HandleClient(c *websocket.Conn) {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		log.Println("ws: missing user_id in locals")
		return
	}

	client := NewClient(h.hub, c, userID, RoleClient)
	h.hub.Register(client)

	go client.WritePump()
	client.ReadPump()
}

func (h *Handler) HandleVet(c *websocket.Conn) {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		log.Println("ws: missing user_id in locals")
		return
	}

	client := NewClient(h.hub, c, userID, RoleVet)
	h.hub.Register(client)

	go client.WritePump()
	client.ReadPump()
}

func AuthMiddlewareForWS(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token query param"})
		}

		claims, err := domain.ParseToken(token, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)

		return c.Next()
	}
}
