// internal/api/router.go
package api

import (
	"budgex_backend/internal/api/handlers"
	"budgex_backend/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Build(db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{AppName: "budgex-backend", DisableStartupMessage: true})

	// public
	api := app.Group("/api")
	handlers.HealthHandler{DB: db}.Register(api)

	// Protected routes
	protected := api.Group("", middleware.FiberAuth())
	handlers.MeHandler{}.Register(protected)
	handlers.TxHandler{DB: db}.Register(protected)
	handlers.CategoryHandler{DB: db}.Register(protected)
	handlers.BudgetHandler{DB: db}.Register(protected)

	return app
}
