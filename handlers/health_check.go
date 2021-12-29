package handlers

import "github.com/gofiber/fiber/v2"

// @Summary      health check
// @Tags         health check
// @Accept       text/html
// @Produce      text/html
// @Success      200  {string} string "OK"
// @Router       /api/ethereum-auth/health [get]
func HealthCheckHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}
