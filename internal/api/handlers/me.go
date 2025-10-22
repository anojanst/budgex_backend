package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type MeHandler struct{}

func (h MeHandler) Register(r fiber.Router) {
	r.Get("/me", h.Me)
}

// Me godoc
// @Summary      Current user id
// @Tags         auth
// @Security     BearerAuth
// @Success      200  {object}  map[string]string  "user_id"
// @Failure      401  {object}  map[string]string  "unauthorized"
// @Router       /me [get]
func (h MeHandler) Me(c *fiber.Ctx) error {
	uid, _ := c.Locals("user_id").(string)
	fmt.Printf("DEBUG: Me handler, user_id: %s\n", uid)
	if uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.JSON(fiber.Map{"user_id": uid})
}
