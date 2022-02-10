package handlers

import "github.com/gofiber/fiber/v2"

// @Summary      metadata
// @Tags         metadata
// @Accept       json
// @Produce      json
// @Success      200  {object} interface{} "metadata json ex: {"signin-text-template": "Sign-in nonce: %s", "ttl-seconds": 5}"
// @Router       /api/ethereum-auth/v1/metadata [get]
func MetadataHandler(signInTextTemplate string, nonceTTL int) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{
			"signin-text-template": signInTextTemplate,
			"ttl-seconds":          nonceTTL,
		})
	}
}
