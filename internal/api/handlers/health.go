package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	DB *gorm.DB
}

func (h HealthHandler) Register(r fiber.Router) {
	r.Get("/health", h.health)
}

func (h HealthHandler) health(c *fiber.Ctx) error {
	// Basic DB ping; if DB is down, still return 503 instead of panicking
	if h.DB != nil {
		if err := h.DB.Exec("SELECT 1").Error; err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"ok":     false,
				"error":  "db_unavailable",
				"detail": err.Error(),
			})
		}
	}
	return c.JSON(fiber.Map{"ok": true, "service": "budgex", "version": "0.1.0"})
}
