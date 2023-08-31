package routes

import (
	handlers "server/handlers/admin/food"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Food_Routes(app *fiber.App) {
	r := app.Group("/api/admin/food")

	// r.Get("/", func(c *fiber.Ctx) error { return handlers.Get_Indredients(c, setup.Admin_DB) })
	r.Post("/details", func(c *fiber.Ctx) error { return handlers.Post_Food_Details(c, setup.Admin_DB) })
	r.Post("/images/req", func(c *fiber.Ctx) error { return handlers.Post_Images_Req(c, setup.Admin_DB) })
	r.Post("/images/confirm", func(c *fiber.Ctx) error { return handlers.Post_Images_Confirm(c, setup.Admin_DB) })
	// r.Post("/thumbnail/req", func(c *fiber.Ctx) error { return handlers.Post_Thumbnail_Req(c, setup.Admin_DB) })
	// r.Post("/thumbnail/confirm", func(c *fiber.Ctx) error { return handlers.Post_Thumbnail_Confirm(c, setup.Admin_DB) })
}
