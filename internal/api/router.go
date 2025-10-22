// internal/api/router.go
package api

import (
	"budgex_backend/internal/api/handlers"
	"budgex_backend/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gorm.io/gorm"
)

func Build(db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{AppName: "budgex-backend", DisableStartupMessage: true})

	// Order matters: requestid -> otel -> auth -> logz
	app.Use(requestid.New())   // adds c.Locals("requestid")
	app.Use(middleware.OTel()) // starts OTel spans (trace/parent propagation)

	// Public
	api := app.Group("/api")
	handlers.HealthHandler{DB: db}.Register(api)

	// Auth (your working Clerk middleware)
	protected := api.Group("", middleware.FiberAuth())

	// Structured logging AFTER auth so user_id is set for logs
	protected.Use(middleware.Logz())
	// Protected routes
	handlers.MeHandler{}.Register(protected)
	handlers.TxHandler{DB: db}.Register(protected)
	handlers.CategoryHandler{DB: db}.Register(protected)
	handlers.BudgetHandler{DB: db}.Register(protected)

	return app
}
