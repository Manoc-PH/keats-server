package routes

import (
	handlers "server/handlers/tracker"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Tracker_Routes(app *fiber.App) {
	r := app.Group("/api/tracker")

	r.Get("/daily_nutrients", func(c *fiber.Ctx) error { return handlers.Get_Daily_Nutrients(c, setup.DB) })
	r.Get("/daily_nutrients_list", func(c *fiber.Ctx) error { return handlers.Get_Daily_Nutrients_List(c, setup.DB) })
	r.Get("/intakes", func(c *fiber.Ctx) error { return handlers.Get_Intakes(c, setup.DB) })
	r.Get("/intake", func(c *fiber.Ctx) error { return handlers.Get_Intake_Details(c, setup.DB) })
	r.Post("/intake", func(c *fiber.Ctx) error { return handlers.Post_Intake(c, setup.DB) })
	r.Put("/intake", func(c *fiber.Ctx) error { return handlers.Put_Intake(c, setup.DB) })
	r.Delete("/intake", func(c *fiber.Ctx) error { return handlers.Delete_Intake(c, setup.DB) })
}
