package routes

import (
	handlers "server/handlers/consumer/common"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Common_Routes(app *fiber.App) {
	r := app.Group("/api/common")

	r.Get("/activity_levels", func(c *fiber.Ctx) error { return handlers.Get_Activity_Levels(c, setup.DB) })
	r.Get("/diet_plans", func(c *fiber.Ctx) error { return handlers.Get_Diet_Plans(c, setup.DB) })
}
