package routes

import (
	handlers "server/handlers/admin/ingredient"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Ingredient_Routes(app *fiber.App) {
	r := app.Group("/api/admin/ingredient")

	r.Get("/", func(c *fiber.Ctx) error { return handlers.Get_Indredients(c, setup.Admin_DB) })
}
