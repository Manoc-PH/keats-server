package routes

import (
	handlers "kryptoverse-api/handlers/auth"
	"kryptoverse-api/setup"

	"github.com/gofiber/fiber/v2"
)

func Auth_Routes(app *fiber.App) {
	r := app.Group("/api")

	r.Post("/signup", func(c *fiber.Ctx) error { return handlers.Sign_Up(c, setup.DB) })
	r.Post("/login", func(c *fiber.Ctx) error { return handlers.Login(c, setup.DB) })
	r.Get("/user", func(c *fiber.Ctx) error { return handlers.User(c, setup.DB) })
	r.Post("/logout", func(c *fiber.Ctx) error { return handlers.Logout(c, setup.DB) })
	r.Get("/verifytoken", func(c *fiber.Ctx) error { return handlers.Verify_Token(c, setup.DB) })
}
