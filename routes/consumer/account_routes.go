package routes

import (
	handlers "server/handlers/consumer/account"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Account_Routes(app *fiber.App) {
	r := app.Group("/api/account")

	// TODO Create update account vitals handler
	r.Get("/account_vitals", func(c *fiber.Ctx) error { return handlers.Get_Account_Vitals(c, setup.DB) })
	r.Put("/account_vitals", func(c *fiber.Ctx) error { return handlers.Update_Account_Vitals(c, setup.DB) })
}
