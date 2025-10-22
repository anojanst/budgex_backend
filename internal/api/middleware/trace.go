package middleware

import (
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
)

// Creates server spans per request and propagates context.
// Trace/Span IDs will be available to exporters and can be read in handlers via otel/trace.
func OTel() fiber.Handler {
	return otelfiber.Middleware() // default: operation name = route
}
