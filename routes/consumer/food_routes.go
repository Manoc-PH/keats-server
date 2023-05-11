package routes

import (
	handlers "server/handlers/consumer/food"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Food_Routes(app *fiber.App) {
	r := app.Group("/api/food")

	r.Get("/search_food", func(c *fiber.Ctx) error { return handlers.Get_Search_Food(c, setup.DB) })
	r.Get("", func(c *fiber.Ctx) error { return handlers.Get_Food_Details(c, setup.DB) })
}
