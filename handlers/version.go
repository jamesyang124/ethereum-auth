package handlers

import "github.com/gofiber/fiber/v2"

// @Summary      version
// @Tags         version
// @Accept       text/html
// @Produce      text/html
// @Success      200  {string} string "respond current service version ex: 1.0.1"
// @Router       /version [get]
func VersionHandler(appVersion string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString(appVersion)
	}
}
