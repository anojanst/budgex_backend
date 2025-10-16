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

	// Placeholders for next tasks:
	// transactions.Register(api, db)
	// categories.Register(api, db)
	// budgets.Register(api, db)

	return app
}
