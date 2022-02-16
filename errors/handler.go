package errors

import "github.com/gofiber/fiber/v2"

func ErrorResponseHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		if e.Code != fiber.StatusFailedDependency && e.Code != fiber.StatusInternalServerError {
			code = fiber.StatusBadRequest
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		} else {
			code = e.Code
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		}
	}

	if code == fiber.StatusBadRequest {
		return c.Status(code).JSON(err)
	}

	return c.Status(code).SendString(err.Error())
}
