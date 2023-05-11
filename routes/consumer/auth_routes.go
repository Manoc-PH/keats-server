package routes

import (
	handlers "server/handlers/consumer/auth"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Auth_Routes(app *fiber.App) {
	r := app.Group("/api/auth")

	r.Post("/signup", func(c *fiber.Ctx) error { return handlers.Sign_Up(c, setup.DB) })
	r.Post("/login", func(c *fiber.Ctx) error { return handlers.Login(c, setup.DB) })
	r.Post("/logout", func(c *fiber.Ctx) error { return handlers.Logout(c, setup.DB) })
}
