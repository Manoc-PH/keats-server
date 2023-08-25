package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/admin/ingredient"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Put_Ingredient_Details(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Put_Ingredient_Details | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	// admin validation
	isAdmin := middlewares.IsAdmin(owner_id, db)
	if isAdmin != true {
		log.Println("Put_Ingredient_Details | Error on auth middleware (Not Admin): ")
		return utilities.Send_Error(c, "Only admin users are allowed to access this endpoint", fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Put_Ingredient_Details)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Put_Ingredient_Details | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// querying ingredient
	old_ingredient_details := schemas.Ingredient_Details{}
	row := query_ingredient_mapping(db, reqData.Ingredient_Mapping_ID)
	if err = scan_ingredient_mapping(row, &old_ingredient_details); err != nil {
		log.Println("Put_Ingredient_Details | Error on scanning ingredient: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON("Could not find ingredient")
	}

	// updating ingredient
	update_ingredient_details(&old_ingredient_details, &reqData.Ingredient_Details)

	// TODO FINISH THIS ENDPOINT
	return c.Status(fiber.StatusOK).JSON(old_ingredient_details)
}

func query_ingredient_mapping(db *sql.DB, ingredient_mapping_id uint) *sql.Row {
	row := db.QueryRow(`SELECT
			ingredient.id, ingredient.name, coalesce(ingredient.name_ph, ''), ingredient.name_owner, coalesce(ingredient.ingredient_desc, ''), category_id,
			ingredient_variant.id, ingredient_variant.name, coalesce(ingredient_variant.name_ph, ''), 
			ingredient_subvariant.id, ingredient_subvariant.name, coalesce(ingredient_subvariant.name_ph, ''),
			nutrient.amount_unit, nutrient.amount_unit_desc, nutrient.serving_size, nutrient.serving_total
		FROM ingredient_mapping
		JOIN ingredient ON ingredient_mapping.ingredient_id = ingredient.id
		JOIN ingredient_variant ON ingredient_mapping.ingredient_variant_id = ingredient_variant.id
		JOIN ingredient_subvariant ON ingredient_mapping.ingredient_subvariant_id = ingredient_subvariant.id
		JOIN nutrient ON ingredient_mapping.nutrient_id = nutrient.id
		WHERE ingredient_mapping.id = $1`,
		ingredient_mapping_id,
	)
	return row
}
func scan_ingredient_mapping(row *sql.Row, ingredient_mapping *schemas.Ingredient_Details) error {
	if err := row.
		Scan(
			&ingredient_mapping.Ingredient.ID,
			&ingredient_mapping.Ingredient.Name,
			&ingredient_mapping.Ingredient.Name_Ph,
			&ingredient_mapping.Ingredient.Name_Owner,
			&ingredient_mapping.Ingredient.Ingredient_Desc,
			&ingredient_mapping.Ingredient.Category_Id,

			&ingredient_mapping.Ingredient_Variant.ID,
			&ingredient_mapping.Ingredient_Variant.Name,
			&ingredient_mapping.Ingredient_Variant.Name_Ph,

			&ingredient_mapping.Ingredient_Subvariant.ID,
			&ingredient_mapping.Ingredient_Subvariant.Name,
			&ingredient_mapping.Ingredient_Subvariant.Name_Ph,

			&ingredient_mapping.Nutrient.Amount_Unit,
			&ingredient_mapping.Nutrient.Amount_Unit_Desc,
			&ingredient_mapping.Nutrient.Serving_Size,
			&ingredient_mapping.Nutrient.Serving_Total,
		); err != nil {
		return err
	}
	return nil
}
func update_ingredient_details(old_ingredient *schemas.Ingredient_Details, new_ingredient *schemas.Ingredient_Details) {
	// Ingredient
	if new_ingredient.Ingredient.Name != "" {
		old_ingredient.Ingredient.Name = new_ingredient.Ingredient.Name
	}
	if new_ingredient.Ingredient.Name_Ph != "" {
		old_ingredient.Ingredient.Name_Ph = new_ingredient.Ingredient.Name_Ph
	}
	if new_ingredient.Ingredient.Name_Owner != "" {
		old_ingredient.Ingredient.Name_Owner = new_ingredient.Ingredient.Name_Owner
	}
	if new_ingredient.Ingredient.Ingredient_Desc != "" {
		old_ingredient.Ingredient.Ingredient_Desc = new_ingredient.Ingredient.Ingredient_Desc
	}
	if new_ingredient.Ingredient.Category_Id != 0 {
		old_ingredient.Ingredient.Category_Id = new_ingredient.Ingredient.Category_Id
	}

	// Variant
	if new_ingredient.Ingredient_Variant.Name != "" {
		old_ingredient.Ingredient_Variant.Name = new_ingredient.Ingredient_Variant.Name
	}
	if new_ingredient.Ingredient_Variant.Name_Ph != "" {
		old_ingredient.Ingredient_Variant.Name_Ph = new_ingredient.Ingredient_Variant.Name_Ph
	}

	// Subvariant
	if new_ingredient.Ingredient_Subvariant.Name != "" {
		old_ingredient.Ingredient_Subvariant.Name = new_ingredient.Ingredient_Subvariant.Name
	}
	if new_ingredient.Ingredient_Subvariant.Name_Ph != "" {
		old_ingredient.Ingredient_Subvariant.Name_Ph = new_ingredient.Ingredient_Subvariant.Name_Ph
	}

	// Nutrient
	if new_ingredient.Nutrient.Amount_Unit != "" {
		old_ingredient.Nutrient.Amount_Unit = new_ingredient.Nutrient.Amount_Unit
	}
	if new_ingredient.Nutrient.Amount_Unit_Desc != "" {
		old_ingredient.Nutrient.Amount_Unit_Desc = new_ingredient.Nutrient.Amount_Unit_Desc
	}
	if new_ingredient.Nutrient.Serving_Size != 0 {
		old_ingredient.Nutrient.Serving_Size = new_ingredient.Nutrient.Serving_Size
	}
	if new_ingredient.Nutrient.Serving_Total != 0 {
		old_ingredient.Nutrient.Serving_Total = new_ingredient.Nutrient.Serving_Total
	}
}
