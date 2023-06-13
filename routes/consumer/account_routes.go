package routes

import (
	handlers "server/handlers/consumer/account"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Account_Routes(app *fiber.App) {
	r := app.Group("/api/account")

	// TODO Create update account vitals handler
	r.Get("/consumer_vitals", func(c *fiber.Ctx) error { return handlers.Get_Consumer_Vitals(c, setup.DB) })
	r.Put("/consumer_vitals", func(c *fiber.Ctx) error { return handlers.Update_Consumer_Vitals(c, setup.DB) })
}
