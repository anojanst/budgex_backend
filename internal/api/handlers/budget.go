package handlers

import (
	"time"

	"budgex_backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BudgetHandler struct{ DB *gorm.DB }

func (h BudgetHandler) Register(r fiber.Router) {
	grp := r.Group("/budgets")
	grp.Get("/", h.List)    // ?month=YYYY-MM (optional; defaults to current)
	grp.Post("/", h.Upsert) // set/replace a budget row
}

type upsertBudgetDTO struct {
	Month      string  `json:"month"`       // "YYYY-MM"
	CategoryID string  `json:"category_id"` // required
	Amount     float64 `json:"amount"`      // required
}

func (h BudgetHandler) List(c *fiber.Ctx) error {
	month := c.Query("month")
	if month == "" {
		now := time.Now().UTC()
		month = now.Format("2006-01")
	}
	var out []models.Budget
	if err := h.DB.Where("user_id = ? AND month = ? AND deleted_at IS NULL",
		userID(c), month).Order("category_id").Find(&out).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(out)
}

func (h BudgetHandler) Upsert(c *fiber.Ctx) error {
	var in upsertBudgetDTO
	if err := c.BodyParser(&in); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad_json"})
	}
	if in.Month == "" || len(in.Month) != 7 {
		return c.Status(422).JSON(fiber.Map{"error": "month_format_YYYY-MM"})
	}
	if in.CategoryID == "" || in.Amount < 0 {
		return c.Status(422).JSON(fiber.Map{"error": "category_id_and_amount_required"})
	}

	row := models.Budget{
		Base:       models.Base{UserID: userID(c)},
		Month:      in.Month,
		CategoryID: in.CategoryID,
		Amount:     in.Amount,
	}
	// Unique on (user_id, month, category_id) effectively
	// GORM upsert with ON CONFLICT
	if err := h.DB.
		Clauses(
			// import clause: "gorm.io/gorm/clause"
			clauseOnConflictSetAmount(),
		).
		Create(&row).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(row)
}
