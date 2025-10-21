package handlers

import (
	"budgex_backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TxHandler struct{ DB *gorm.DB }

func (h TxHandler) Register(r fiber.Router) {
	tx := r.Group("/transactions")
	tx.Get("/", h.List)
	tx.Post("/", h.Create)
}

type createTxDTO struct {
	Type       string  `json:"type"`           // "income" | "expense"
	Date       *string `json:"date,omitempty"` // ISO; defaults now
	Amount     float64 `json:"amount"`
	Payee      *string `json:"payee"`
	Memo       *string `json:"memo"`
	CategoryID *string `json:"category_id"`
	Tags       *string `json:"tags"`
}

func userID(c *fiber.Ctx) string {
	if v := c.Locals("user_id"); v != nil {
		if s, _ := v.(string); s != "" {
			return s
		}
	}
	return ""
}

func (h TxHandler) List(c *fiber.Ctx) error {
	var out []models.Transaction
	err := h.DB.
		Where("user_id = ? AND deleted_at IS NULL", userID(c)).
		Order("date DESC, created_at DESC").
		Limit(100).
		Find(&out).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(out)
}

func (h TxHandler) Create(c *fiber.Ctx) error {
	var in createTxDTO
	if err := c.BodyParser(&in); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad_json"})
	}
	if in.Type != "income" && in.Type != "expense" {
		return c.Status(422).JSON(fiber.Map{"error": "type_must_be_income_or_expense"})
	}
	d := time.Now().UTC()
	if in.Date != nil && *in.Date != "" {
		if t, err := time.Parse(time.RFC3339, *in.Date); err == nil {
			d = t
		}
	}
	tx := models.Transaction{
		Base: models.Base{UserID: userID(c)},
		Type: in.Type, Date: d, Amount: in.Amount,
		Payee: in.Payee, Memo: in.Memo, CategoryID: in.CategoryID, Tags: in.Tags,
	}
	if err := h.DB.Create(&tx).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(tx)
}
