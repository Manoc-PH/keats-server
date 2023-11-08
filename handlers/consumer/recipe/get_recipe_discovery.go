package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Get_Recipe_Discovery(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Recipe_Discovery | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	response := new(schemas.Res_Get_Recipe_Discovery)
	// * getting recipe details
	err = get_recipe_discovery(db, &response.Recipes)
	if err != nil {
		log.Println("Get_Recipe_Discovery | Error on get_recipe_discovery: ", err.Error())
		return utilities.Send_Error(c, "An error occured in fetching recipe", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func get_recipe_discovery(db *sql.DB, recipe *[]models.Recipe) error {
	local_recipes := []models.Recipe{}
	rows, err := db.Query(`
		SELECT 
			id,
			name,
			name_ph,
			name_owner,
			image_url,
			image_name,
			rating,
			rating_count
		FROM recipe
		WHERE date_created >= CURRENT_DATE - INTERVAL '31 days'
		ORDER BY rating DESC, date_created DESC
		LIMIT 10`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var local_rec = models.Recipe{}
		if err := rows.
			Scan(
				&local_rec.ID,
				&local_rec.Name,
				&local_rec.Name_Ph,
				&local_rec.Name_Owner,
				&local_rec.Image_URL,
				&local_rec.Image_Name,
				&local_rec.Rating,
				&local_rec.Rating_Count,
			); err != nil {
			return err
		}
		local_recipes = append(local_recipes, local_rec)
	}
	if len(local_recipes) > 5 {
		recipe = &local_recipes
		return nil
	}
	return nil
}
