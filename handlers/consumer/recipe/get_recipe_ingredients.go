package handlers

import (
	"database/sql"
	"log"
	"server/constants"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Get_Recipe_Ingredients(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("get_recipe_ingredients | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Recipe_Details)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("get_recipe_ingredients | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	response := new(schemas.Res_Get_Recipe_Ingredients)
	// * getting recipe details
	err = get_recipe_ingredients(db, reqData.Recipe_Id, &response.Ingredients)
	if err != nil {
		log.Println("get_recipe_ingredients | Error on Get_Recipe_Ingredients: ", err.Error())
		return utilities.Send_Error(c, "An error occured in fetching recipe", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func get_recipe_ingredients(db *sql.DB, recipe_id uuid.UUID, recipe_ings *[]schemas.Recipe_Ingredient_Details_Schema) error {
	rows, err := db.Query(`SELECT
			recipe_ingredient.id,
			recipe_ingredient.ingredient_mapping_id,
			recipe_ingredient.food_id,
			recipe_ingredient.amount,
			recipe_ingredient.amount_unit,
			recipe_ingredient.amount_unit_desc,
			recipe_ingredient.serving_size,
			COALESCE(food.name, ''),
			COALESCE(food.name_owner, ''),
			COALESCE(ingredient.name, ''),
			COALESCE(ingredient_variant.name, ''),
			COALESCE(ingredient_subvariant.name, ''),
			COALESCE(ingredient.name_owner, '') 
		FROM recipe_ingredient
		LEFT JOIN food on recipe_ingredient.food_id = food.id
		LEFT JOIN ingredient_mapping ON recipe_ingredient.ingredient_mapping_id = ingredient_mapping.id
		LEFT JOIN ingredient ON ingredient_mapping.ingredient_id = ingredient.id
		LEFT JOIN ingredient_variant ON ingredient_mapping.ingredient_variant_id = ingredient_variant.id
		LEFT JOIN ingredient_subvariant ON ingredient_mapping.ingredient_subvariant_id = ingredient_subvariant.id
		WHERE recipe_ingredient.recipe_id = $1`,
		recipe_id,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var recipe_ing = schemas.Recipe_Ingredient_Details_Schema{}
		var food_name string
		var food_owner_name string
		var ingredient_name string
		var ingredient_variant_name string
		var ingredient_subvariant_name string
		var ingredient_owner_name string
		if err := rows.
			Scan(
				&recipe_ing.ID,
				&recipe_ing.Ingredient_Mapping_Id,
				&recipe_ing.Food_Id,
				&recipe_ing.Amount,
				&recipe_ing.Amount_Unit,
				&recipe_ing.Amount_Unit_Desc,
				&recipe_ing.Serving_Size,
				&food_name,
				&food_owner_name,
				&ingredient_name,
				&ingredient_variant_name,
				&ingredient_subvariant_name,
				&ingredient_owner_name,
			); err != nil {
			return err
		}
		if recipe_ing.Ingredient_Mapping_Id != constants.Empty_UUID {
			recipe_ing.Name = ingredient_name + " " + ingredient_variant_name + " " + ingredient_subvariant_name
			recipe_ing.Name_Owner = ingredient_owner_name
		}
		if recipe_ing.Food_Id != constants.Empty_UUID {
			recipe_ing.Name = food_name
			recipe_ing.Name_Owner = food_owner_name
		}
		*recipe_ings = append(*recipe_ings, recipe_ing)
	}
	return nil
}
