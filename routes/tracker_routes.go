package routes

import (
	handlers "server/handlers/tracker"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Tracker_Routes(app *fiber.App) {
	r := app.Group("/api/tracker")

	r.Get("/macros", func(c *fiber.Ctx) error { return handlers.Get_Macros(c, setup.DB) })
	r.Get("/macros_list", func(c *fiber.Ctx) error { return handlers.Get_Macros_List(c, setup.DB) })
	r.Post("/intake", func(c *fiber.Ctx) error { return handlers.Post_Intake(c, setup.DB) })
}
