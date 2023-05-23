package routes

import (
	handlers "server/handlers/consumer/ingredient"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Ingredient_Routes(app *fiber.App) {
	r := app.Group("/api/ingredient")

	r.Get("/mapping", func(c *fiber.Ctx) error { return handlers.Get_Ingredient_Mapping_Details(c, setup.DB) })
	r.Get("/search_ingredient", func(c *fiber.Ctx) error { return handlers.Get_Search_Ingredient(c, setup.DB_Search) })
}
