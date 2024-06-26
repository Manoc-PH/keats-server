package routes

import (
	handlers "server/handlers/admin/ingredient"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Ingredient_Routes(app *fiber.App) {
	r := app.Group("/api/admin/ingredient")

	r.Get("/", func(c *fiber.Ctx) error { return handlers.Get_Indredients(c, setup.Admin_DB) })
	r.Put("/details", func(c *fiber.Ctx) error { return handlers.Put_Ingredient_Details(c, setup.Admin_DB) })
	r.Delete("/images", func(c *fiber.Ctx) error { return handlers.Delete_Images(c, setup.Admin_DB) })
	r.Post("/images/req", func(c *fiber.Ctx) error { return handlers.Post_Images_Req(c, setup.Admin_DB) })
	r.Post("/images/confirm", func(c *fiber.Ctx) error { return handlers.Post_Images_Confirm(c, setup.Admin_DB) })
	r.Post("/images/duplicate", func(c *fiber.Ctx) error { return handlers.Post_Images_Duplicate(c, setup.Admin_DB) })
	r.Post("/thumbnail/req", func(c *fiber.Ctx) error { return handlers.Post_Thumbnail_Req(c, setup.Admin_DB) })
	r.Post("/thumbnail/confirm", func(c *fiber.Ctx) error { return handlers.Post_Thumbnail_Confirm(c, setup.Admin_DB) })
}
