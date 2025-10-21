// internal/api/middleware/clerk_fiber.go
package middleware

import (
	"net/http"
	"net/http/httptest"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func FiberAuth() fiber.Handler {
	httpMW := clerkhttp.WithHeaderAuthorization() // verifies Bearer & adds claims to req.Context()

	return func(c *fiber.Ctx) error {
		// 1) Convert Fiber ctx -> *http.Request
		req, err := adaptor.ConvertRequest(c, true)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "convert_request_failed"})
		}

		// 2) Build a terminal handler that will run AFTER Clerk middleware
		var userID string
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if sc, ok := clerk.SessionClaimsFromContext(r.Context()); ok && sc != nil {
				userID = sc.Subject
			}
			// do not write a body; we only need claims
		})

		// 3) Run Clerk middleware + next on this request
		rr := httptest.NewRecorder()
		handler := httpMW(next)
		handler.ServeHTTP(rr, req)

		// 4) Check result and expose user_id to Fiber
		if rr.Code >= 400 || userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		c.Locals("user_id", userID)
		return c.Next()
	}
}
