package handlers

import (
	"database/sql"
	"log"
	"server/constants"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func Get_Recipe_Filtered(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("get_recipe_reviews | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Recipe_Filtered)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("get_recipe_ingredients | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	response := new(schemas.Res_Get_Recipe_Filtered)
	// * getting recipe details
	recipes, err := get_recipe_filtered(
		db,
		constants.Recipe_Nutrition_Categories_SQL[reqData.Filter],
		constants.Recipe_Nutrition_Categories_Order[reqData.Filter],
	)
	if err != nil {
		log.Println("get_recipe_reviews | Error on Get_Recipe_Ingredients: ", err.Error())
		return utilities.Send_Error(c, "An error occured in fetching recipe", fiber.StatusInternalServerError)
	}
	response.Recipes = recipes
	return c.Status(fiber.StatusOK).JSON(response)
}

// TODO RETURN THE MACRO NUTRIENTS AND CALORIES AND SHOW IT PER 100g
func get_recipe_filtered(db *sql.DB, filter, order string) ([]models.Recipe, error) {
	local_recipes := []models.Recipe{}
	query := ``
	if order == constants.Recipe_Nutrition_Category_Order.DESC {
		query = `
			SELECT 
				recipe.id,
				recipe.name,
				recipe.name_ph,
				recipe.name_owner,
				recipe.main_image_link,
				recipe.rating,
				recipe.rating_count
			FROM recipe
			JOIN nutrient ON recipe.nutrient_id = nutrient.id
			ORDER BY nutrient.` + pq.QuoteIdentifier(filter) + ` DESC
			LIMIT 10`
	}
	if order == constants.Recipe_Nutrition_Category_Order.ASC {
		query = `
			SELECT 
				recipe.id,
				recipe.name,
				recipe.name_ph,
				recipe.name_owner,
				recipe.main_image_link,
				recipe.rating,
				recipe.rating_count
			FROM recipe
			JOIN nutrient ON recipe.nutrient_id = nutrient.id
			ORDER BY nutrient.` + pq.QuoteIdentifier(filter) + ` ASC
			LIMIT 10`
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
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
				&local_rec.Main_Image_Link,
				&local_rec.Rating,
				&local_rec.Rating_Count,
			); err != nil {
			return nil, err
		}
		local_recipes = append(local_recipes, local_rec)
	}
	return local_recipes, nil
}
