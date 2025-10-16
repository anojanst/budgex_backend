package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Common() []fiber.Handler {
	return []fiber.Handler{
		requestid.New(),
		logger.New(logger.Config{
			// concise dev logger; adjust format if you want JSON logs
			Format: "[${time}] ${ip} ${status} ${latency} ${method} ${path} rid=${locals:requestid}\n",
		}),
	}
}
