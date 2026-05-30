package vets

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/danielm/app_sara_backend/internal/domain"
	"github.com/danielm/app_sara_backend/internal/middleware"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *fiber.Ctx) error {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.Query("radius")

	var lat, lng *float64
	var radius *int

	if latStr != "" && lngStr != "" && radiusStr != "" {
		latF, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid lat"})
		}
		lngF, err := strconv.ParseFloat(lngStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid lng"})
		}
		radiusI, err := strconv.Atoi(radiusStr)
		if err != nil || radiusI <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid radius"})
		}
		lat = &latF
		lng = &lngF
		radius = &radiusI
	}

	vets, err := h.service.List(lat, lng, radius)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(vets)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	vet, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if vet == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "vet not found"})
	}
	return c.JSON(vet)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	userRole := c.Locals("user_role").(domain.UserRole)
	isAdmin := userRole == domain.RoleAdmin

	var input domain.UpdateVetInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	vet, err := h.service.Update(id, userID, input, isAdmin)
	if err != nil {
		if err.Error() == "forbidden" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if vet == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "vet not found"})
	}
	return c.JSON(vet)
}

func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	userRole := c.Locals("user_role").(domain.UserRole)
	isAdmin := userRole == domain.RoleAdmin

	var input domain.UpdateVetStatusInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.UpdateStatus(id, userID, input.Status, isAdmin); err != nil {
		if err.Error() == "forbidden" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": input.Status})
}

func (h *Handler) UpdateLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	var input domain.UpdateVetLocationInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.UpdateLocation(id, userID, input); err != nil {
		if err.Error() == "forbidden" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "location updated"})
}

func RegisterRoutes(r fiber.Router, handler *Handler, authMiddleware, adminMiddleware fiber.Handler, middlewares ...fiber.Handler) {
	v := r.Group("/vets")
	for _, m := range middlewares {
		v.Use(m)
	}

	v.Get("/", authMiddleware, handler.List)
	v.Get("/:id", authMiddleware, handler.GetByID)
	v.Put("/:id", authMiddleware, handler.Update)
	v.Put("/:id/status", authMiddleware, handler.UpdateStatus)
	v.Put("/:id/location", authMiddleware, handler.UpdateLocation)
}
