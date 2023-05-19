package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/ingredient"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Get_Ingredient_Details(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Ingredient_Details | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Ingredient_Details)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Ingredient_Details | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	response := schemas.Res_Get_Ingredient_Details{}
	// querying ingredient
	row := query_ingredient(reqData.Ingredient_Mapping_ID, db)
	err = scan_ingredient(row, &response)
	if err != nil && err == sql.ErrNoRows {
		log.Println("Get_Ingredient_Details | error in scanning food: ", err.Error())
		return utilities.Send_Error(c, "Food does not exist", fiber.StatusBadRequest)
	}
	// Server Error
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get_Ingredient_Details | error in scanning food: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	// TODO Query images of ingredient
	// querying ingredient images
	// images, err := query_and_scan_food_images(db, reqData.Food_ID)
	// if err != nil {
	// 	return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	// }
	// response.Food_Images = images
	return c.Status(fiber.StatusOK).JSON(response)
}

func query_ingredient(ingredient_mapping_id uint, db *sql.DB) *sql.Row {
	row := db.QueryRow(`SELECT
			ingredient.id, ingredient.name, coalesce(ingredient.name_ph, ''), ingredient.name_owner,
			ingredient_variant.id, ingredient_variant.name, coalesce(ingredient_variant.name_ph, ''), 
			ingredient_subvariant.id, ingredient_subvariant.name, coalesce(ingredient_subvariant.name_ph, ''), 
			nutrient.id,
			nutrient.amount,
			nutrient.amount_unit,
			nutrient.amount_unit_desc,
			nutrient.serving_size,
			nutrient.calories,
			nutrient.protein,
			nutrient.carbs,
			nutrient.fats,
			nutrient.trans_fat,
			nutrient.saturated_fat,
			nutrient.sugars,
			nutrient.fiber,
			nutrient.sodium,
			nutrient.iron,
			nutrient.calcium
		FROM ingredient_mapping
		JOIN ingredient ON ingredient_mapping.ingredient_id = ingredient.id
		JOIN ingredient_variant ON ingredient_mapping.ingredient_variant_id = ingredient_variant.id
		JOIN ingredient_subvariant ON ingredient_mapping.ingredient_subvariant_id = ingredient_subvariant.id
		JOIN nutrient ON ingredient_mapping.nutrient_id = nutrient.id
		WHERE ingredient_mapping.id = $1`,
		// casting timestamp to date
		ingredient_mapping_id,
	)
	return row
}
func scan_ingredient(row *sql.Row, ingredient_mapping *schemas.Res_Get_Ingredient_Details) error {
	if err := row.
		Scan(
			&ingredient_mapping.Ingredient.ID,
			&ingredient_mapping.Ingredient.Name,
			&ingredient_mapping.Ingredient.Name_Ph,
			&ingredient_mapping.Ingredient.Name_Owner,

			&ingredient_mapping.Ingredient_Variant.ID,
			&ingredient_mapping.Ingredient_Variant.Name,
			&ingredient_mapping.Ingredient_Variant.Name_Ph,

			&ingredient_mapping.Ingredient_Subvariant.ID,
			&ingredient_mapping.Ingredient_Subvariant.Name,
			&ingredient_mapping.Ingredient_Subvariant.Name_Ph,

			&ingredient_mapping.Nutrient.ID,
			&ingredient_mapping.Nutrient.Amount,
			&ingredient_mapping.Nutrient.Amount_Unit,
			&ingredient_mapping.Nutrient.Amount_Unit_Desc,
			&ingredient_mapping.Nutrient.Serving_Size,
			&ingredient_mapping.Nutrient.Calories,
			&ingredient_mapping.Nutrient.Protein,
			&ingredient_mapping.Nutrient.Carbs,
			&ingredient_mapping.Nutrient.Fats,
			&ingredient_mapping.Nutrient.Trans_Fat,
			&ingredient_mapping.Nutrient.Saturated_Fat,
			&ingredient_mapping.Nutrient.Sugars,
			&ingredient_mapping.Nutrient.Fiber,
			&ingredient_mapping.Nutrient.Sodium,
			&ingredient_mapping.Nutrient.Iron,
			&ingredient_mapping.Nutrient.Calcium,
		); err != nil {
		return err
	}
	return nil
}
