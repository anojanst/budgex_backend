package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ---------- Wire-up ----------
type AnalyticsHandler struct{ DB *gorm.DB }

func (h AnalyticsHandler) Register(r fiber.Router) {
	g := r.Group("/analytics")
	g.Get("/spend_summary", h.SpendSummary)
	g.Get("/cashflow_forecast", h.CashflowForecast)
}

// ---------- DTOs ----------
type SpendSummaryRow struct {
	CategoryID *string `json:"category_id,omitempty"`
	Category   *string `json:"category,omitempty"`
	Total      float64 `json:"total"`
}

type SpendSummaryResp struct {
	Month        string            `json:"month"` // YYYY-MM
	TotalExpense float64           `json:"total_expense"`
	TotalIncome  float64           `json:"total_income"`
	ByCategory   []SpendSummaryRow `json:"by_category"`
}

type CashflowPoint struct {
	Month    string  `json:"month"` // YYYY-MM
	Income   float64 `json:"income"`
	Expense  float64 `json:"expense"`
	Net      float64 `json:"net"`
	Forecast bool    `json:"forecast"` // true for projected months
}

type CashflowResp struct {
	WindowMonths int             `json:"window_months"`
	Horizon      int             `json:"horizon"`
	Points       []CashflowPoint `json:"points"`
}

// ---------- Helpers ----------
func monthParamOrNow(c *fiber.Ctx) (string, time.Time) {
	m := c.Query("month")
	if m == "" {
		now := time.Now().UTC()
		return now.Format("2006-01"), time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	// parse YYYY-MM
	t, _ := time.Parse("2006-01", m)
	if t.IsZero() {
		now := time.Now().UTC()
		return now.Format("2006-01"), time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	return m, t
}

// -----------------------------
// @Summary      Spend summary for a month
// @Tags         analytics
// @Security     BearerAuth
// @Produce      json
// @Param        month  query   string  false  "YYYY-MM (defaults to current)"
// @Success      200    {object}  SpendSummaryResp
// @Failure      401    {object}  map[string]string
// @Router       /analytics/spend_summary [get]
func (h AnalyticsHandler) SpendSummary(c *fiber.Ctx) error {
	uid, _ := c.Locals("user_id").(string)
	if uid == "" {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	monthStr, month := monthParamOrNow(c)
	from := month
	to := month.AddDate(0, 1, 0)

	// Totals (income & expense)
	type totalRow struct {
		Type  string  `gorm:"column:type"`
		Total float64 `gorm:"column:total"`
	}
	var totals []totalRow
	if err := h.DB.Raw(`
		SELECT type, COALESCE(SUM(amount),0) AS total
		FROM transactions
		WHERE user_id = ? AND deleted_at IS NULL
		  AND date >= ? AND date < ?
		GROUP BY type
	`, uid, from, to).Scan(&totals).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	var totalIncome, totalExpense float64
	for _, t := range totals {
		if t.Type == "income" {
			totalIncome = t.Total
		} else if t.Type == "expense" {
			totalExpense = t.Total
		}
	}

	// By category (expenses)
	type catRow struct {
		CategoryID *string `gorm:"column:category_id"`
		Category   *string `gorm:"column:name"`
		Total      float64 `gorm:"column:total"`
	}
	var catRows []catRow
	if err := h.DB.Raw(`
		SELECT t.category_id,
		       c.name,
		       COALESCE(SUM(t.amount),0) AS total
		FROM transactions t
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = ? AND t.deleted_at IS NULL
		  AND t.type = 'expense'
		  AND t.date >= ? AND t.date < ?
		GROUP BY t.category_id, c.name
		ORDER BY total DESC NULLS LAST
	`, uid, from, to).Scan(&catRows).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	byCat := make([]SpendSummaryRow, 0, len(catRows))
	for _, r := range catRows {
		byCat = append(byCat, SpendSummaryRow{
			CategoryID: r.CategoryID,
			Category:   r.Category,
			Total:      r.Total,
		})
	}

	return c.JSON(SpendSummaryResp{
		Month:        monthStr,
		TotalExpense: totalExpense,
		TotalIncome:  totalIncome,
		ByCategory:   byCat,
	})
}

// -----------------------------
// @Summary      Cashflow forecast (simple average projection)
// @Tags         analytics
// @Security     BearerAuth
// @Produce      json
// @Param        window_months  query  int  false  "How many past months to average (default 3)" minimum(1) maximum(12)
// @Param        horizon        query  int  false  "How many future months to forecast (default 3)" minimum(1) maximum(12)
// @Success      200    {object}  CashflowResp
// @Failure      401    {object}  map[string]string
// @Router       /analytics/cashflow_forecast [get]
func (h AnalyticsHandler) CashflowForecast(c *fiber.Ctx) error {
	uid, _ := c.Locals("user_id").(string)
	if uid == "" {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	// Params
	win := c.QueryInt("window_months", 3)
	if win < 1 {
		win = 1
	}
	if win > 12 {
		win = 12
	}
	hz := c.QueryInt("horizon", 3)
	if hz < 1 {
		hz = 1
	}
	if hz > 12 {
		hz = 12
	}

	// Past N months aggregates
	type row struct {
		Month   string  `gorm:"column:month"`
		Income  float64 `gorm:"column:income"`
		Expense float64 `gorm:"column:expense"`
		Net     float64 `gorm:"column:net"`
	}
	var past []row
	if err := h.DB.Raw(`
		WITH m AS (
		  SELECT to_char(date_trunc('month', date), 'YYYY-MM') AS month,
		         SUM(CASE WHEN type='income'  THEN amount ELSE 0 END) AS income,
		         SUM(CASE WHEN type='expense' THEN amount ELSE 0 END) AS expense
		  FROM transactions
		  WHERE user_id = ? AND deleted_at IS NULL
		  GROUP BY 1
		)
		SELECT month, income, expense, (income - expense) AS net
		FROM m
		ORDER BY month DESC
		LIMIT ?
	`, uid, win).Scan(&past).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	points := []CashflowPoint{}
	var avgIncome, avgExpense float64
	if len(past) > 0 {
		// reverse to chronological order
		for i := len(past) - 1; i >= 0; i-- {
			points = append(points, CashflowPoint{
				Month:    past[i].Month,
				Income:   past[i].Income,
				Expense:  past[i].Expense,
				Net:      past[i].Net,
				Forecast: false,
			})
			avgIncome += past[i].Income
			avgExpense += past[i].Expense
		}
		avgIncome /= float64(len(past))
		avgExpense /= float64(len(past))
	}

	// Forecast next H months using simple average (can be replaced with ARIMA later)
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= hz; i++ {
		fm := start.AddDate(0, i, 0)
		month := fm.Format("2006-01")
		points = append(points, CashflowPoint{
			Month:    month,
			Income:   avgIncome,
			Expense:  avgExpense,
			Net:      avgIncome - avgExpense,
			Forecast: true,
		})
	}

	return c.JSON(CashflowResp{
		WindowMonths: win,
		Horizon:      hz,
		Points:       points,
	})
}
