package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/session/v2"
)

var SessionStore *session.Session

func FiberMiddleware(a *fiber.App) {
	a.Use(
		// Add CORS to each routes
		cors.New(),
		// simple logger
		logger.New(),
	)

	SessionStore = session.New(session.Config{
		Lookup:   "cookie:sessionID",
		Secure:   true,
		SameSite: "Lax",
	})

	a.Use(SessionStore)
}
