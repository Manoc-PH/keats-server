package routes

import (
	handlers "server/handlers/consumer/recipe"
	"server/setup"

	"github.com/gofiber/fiber/v2"
)

func Recipe_Routes(app *fiber.App) {
	r := app.Group("/api/recipe")

	r.Post("", func(c *fiber.Ctx) error { return handlers.Post_Recipe(c, setup.DB, setup.DB_Search) })
	r.Post("/images", func(c *fiber.Ctx) error { return handlers.Post_Images(c, setup.DB) })
	r.Post("/review", func(c *fiber.Ctx) error { return handlers.Post_Recipe_Review(c, setup.DB, setup.DB_Search) })
	r.Post("/like", func(c *fiber.Ctx) error { return handlers.Post_Recipe_Like(c, setup.DB) })
	r.Get("", func(c *fiber.Ctx) error { return handlers.Get_Recipe_Details(c, setup.DB) })
	r.Get("/ingredients", func(c *fiber.Ctx) error { return handlers.Get_Recipe_Ingredients(c, setup.DB) })
	r.Get("/instructions", func(c *fiber.Ctx) error { return handlers.Get_Recipe_Instructions(c, setup.DB) })
	r.Get("/reviews", func(c *fiber.Ctx) error { return handlers.Get_Recipe_Reviews(c, setup.DB) })
	r.Get("/review", func(c *fiber.Ctx) error { return handlers.Get_Recipe_Review(c, setup.DB) })
	r.Get("/actions", func(c *fiber.Ctx) error { return handlers.Get_Recipe_Actions(c, setup.DB) })
	r.Get("/search", func(c *fiber.Ctx) error { return handlers.Get_Search_Recipe(c, setup.DB_Search) })
	r.Get("/discovery", func(c *fiber.Ctx) error { return handlers.Get_Recipe_Discovery(c, setup.DB) })
	r.Get("/filtered", func(c *fiber.Ctx) error { return handlers.Get_Recipe_Filtered(c, setup.DB) })
	r.Patch("", func(c *fiber.Ctx) error { return handlers.Patch_Recipe(c, setup.DB, setup.DB_Search) })
	r.Patch("/review", func(c *fiber.Ctx) error { return handlers.Patch_Recipe_Review(c, setup.DB, setup.DB_Search) })
	r.Delete("", func(c *fiber.Ctx) error { return handlers.Delete_Recipe(c, setup.DB, setup.DB_Search) })
	r.Delete("/like", func(c *fiber.Ctx) error { return handlers.Delete_Recipe_Like(c, setup.DB) })
	r.Delete("/review", func(c *fiber.Ctx) error { return handlers.Delete_Recipe_Review(c, setup.DB) })
}
