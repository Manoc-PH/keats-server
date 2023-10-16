package routes

import (
	handlers "server/handlers/consumer/recipe"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Recipe_Routes(app *fiber.App) {
	r := app.Group("/api/recipe")

	r.Post("", func(c *fiber.Ctx) error { return handlers.Post_Recipe(c, setup.DB) })
	r.Post("/images", func(c *fiber.Ctx) error { return handlers.Post_Images(c, setup.DB) })
}
