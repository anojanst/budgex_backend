package api

import (
	"budgex_backend/internal/api/handlers"
	"budgex_backend/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Build(db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "budgex-backend",
		DisableStartupMessage: true,
	})

	// Global middleware
	for _, mw := range middleware.Common() {
		app.Use(mw)
	}

	// API group
	api := app.Group("/api")

	// Health
	handlers.HealthHandler{DB: db}.Register(api)
	handlers.TxHandler{DB: db}.Register(api)
	handlers.CategoryHandler{DB: db}.Register(api)
	handlers.BudgetHandler{DB: db}.Register(api)

	return app
}
