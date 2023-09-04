package routes

import (
	handlers "server/handlers/admin/auth"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Auth_Routes(app *fiber.App) {
	r := app.Group("/api/admin/auth")

	r.Post("/login", func(c *fiber.Ctx) error { return handlers.Login(c, setup.DB) })
	r.Post("/logout", func(c *fiber.Ctx) error { return handlers.Logout(c, setup.DB) })
}
