package handlers

import "github.com/gofiber/fiber/v2"

// @Summary      metadata
// @Tags         metadata
// @Accept       json
// @Produce      json
// @Success      200  {object} interface{} "metadata json ex: {"signin-text-template": "Sign-in nonce: %s"}"
// @Router       /metadata [get]
func MetadataHandler(signInTextTemplate string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{"signin-text-template": signInTextTemplate})
	}
}
