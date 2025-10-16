package handlers

import (
	"budgex_backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoryHandler struct{ DB *gorm.DB }

func (h CategoryHandler) Register(r fiber.Router) {
	grp := r.Group("/categories")
	grp.Get("/", h.List)
	grp.Post("/", h.Create)
}

type createCategoryDTO struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"`
}

func (h CategoryHandler) List(c *fiber.Ctx) error {
	var out []models.Category
	if err := h.DB.Where("user_id = ? AND deleted_at IS NULL", userID(c)).
		Order("name ASC").Find(&out).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(out)
}

func (h CategoryHandler) Create(c *fiber.Ctx) error {
	var in createCategoryDTO
	if err := c.BodyParser(&in); err != nil || in.Name == "" {
		return c.Status(422).JSON(fiber.Map{"error": "name_required"})
	}
	cat := models.Category{
		Base:     models.Base{UserID: userID(c)},
		Name:     in.Name,
		ParentID: in.ParentID,
	}
	if err := h.DB.Create(&cat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(cat)
}
