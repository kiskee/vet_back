package services

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

func (h *Handler) List(c *fiber.Ctx) error {
	services, err := h.service.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(services)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	svc, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if svc == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "service not found"})
	}
	return c.JSON(svc)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var input domain.CreateServiceInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	svc, err := h.service.Create(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(svc)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var input domain.UpdateServiceInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	svc, err := h.service.Update(id, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if svc == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "service not found"})
	}
	return c.JSON(svc)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func RegisterRoutes(r fiber.Router, handler *Handler, authMiddleware, adminMiddleware fiber.Handler, middlewares ...fiber.Handler) {
	svc := r.Group("/services")
	for _, m := range middlewares {
		svc.Use(m)
	}

	svc.Get("/", authMiddleware, handler.List)
	svc.Get("/:id", authMiddleware, handler.GetByID)
	svc.Post("/", authMiddleware, adminMiddleware, handler.Create)
	svc.Put("/:id", authMiddleware, adminMiddleware, handler.Update)
	svc.Delete("/:id", authMiddleware, adminMiddleware, handler.Delete)
}
