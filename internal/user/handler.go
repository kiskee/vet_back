package user

import (
	"github.com/gofiber/fiber/v2"

	"github.com/danielm/app_sara_backend/internal/domain"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := h.service.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}

func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var input domain.UpdateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	user, err := h.service.UpdateProfile(userID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}

func (h *Handler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.service.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}

func (h *Handler) ListUsers(c *fiber.Ctx) error {
	users, err := h.service.ListUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(users)
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func RegisterRoutes(r fiber.Router, handler *Handler, authMiddleware, adminMiddleware fiber.Handler) {
	users := r.Group("/users")

	users.Get("/me", authMiddleware, handler.GetProfile)
	users.Put("/me", authMiddleware, handler.UpdateProfile)

	users.Get("/", authMiddleware, adminMiddleware, handler.ListUsers)
	users.Get("/:id", authMiddleware, adminMiddleware, handler.GetUser)
	users.Delete("/:id", authMiddleware, adminMiddleware, handler.DeleteUser)
}
