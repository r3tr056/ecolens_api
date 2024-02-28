package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/r3tr056/ecolens_api/platform/db"
)

// healthCheck godoc
// @Summary Perform a health check
// @Description Checks the health status of the application and its dependencies.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Health status response"
// @Router /health [get]
func healthCheck(c *fiber.Ctx) error {
	// health status
	status := map[string]string{
		"status": "ok",
	}

	if db.DatabaseCheck() {
		status["postgres_database"] = "ok"
	} else {
		status["postgres_database"] = "error"
	}

	return c.JSON(status)
}
