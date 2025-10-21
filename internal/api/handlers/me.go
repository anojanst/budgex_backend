package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type MeHandler struct{}

func (h MeHandler) Register(r fiber.Router) {
	r.Get("/me", h.Me)
}

func (h MeHandler) Me(c *fiber.Ctx) error {
	uid, _ := c.Locals("user_id").(string)
	fmt.Printf("DEBUG: Me handler, user_id: %s\n", uid)
	if uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.JSON(fiber.Map{"user_id": uid})
}
