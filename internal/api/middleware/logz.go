package middleware

import (
	"time"

	"budgex_backend/internal/observability"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func Logz() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		lat := time.Since(start)

		// Correlate with OTel trace/span if present
		spanCtx := trace.SpanContextFromContext(c.UserContext())
		traceID := ""
		spanID := ""
		if spanCtx.HasTraceID() {
			traceID = spanCtx.TraceID().String()
		}
		if spanCtx.HasSpanID() {
			spanID = spanCtx.SpanID().String()
		}

		uid, _ := c.Locals("user_id").(string)
		reqID := c.Locals("requestid")
		logger := observability.L()

		fields := []zap.Field{
			zap.Int("status", c.Response().StatusCode()),
			zap.String("method", c.Method()),
			zap.String("path", c.OriginalURL()),
			zap.String("ip", c.IP()),
			zap.Duration("latency_ms", lat),
		}
		if traceID != "" {
			fields = append(fields, zap.String("trace_id", traceID))
		}
		if spanID != "" {
			fields = append(fields, zap.String("span_id", spanID))
		}
		if uid != "" {
			fields = append(fields, zap.String("user_id", uid))
		}
		if rid, ok := reqID.(string); ok && rid != "" {
			fields = append(fields, zap.String("request_id", rid))
		}

		// Log level by status class
		switch {
		case c.Response().StatusCode() >= 500:
			logger.Error("http_request", fields...)
		case c.Response().StatusCode() >= 400:
			logger.Warn("http_request", fields...)
		default:
			logger.Info("http_request", fields...)
		}
		return err
	}
}
