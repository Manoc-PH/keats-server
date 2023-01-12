package utilities

import "github.com/gofiber/fiber/v2"

func Send_Error(c *fiber.Ctx, error_msg string, status int) error {
	return c.Status(status).JSON(fiber.Map{
		"message": error_msg,
	})
}
