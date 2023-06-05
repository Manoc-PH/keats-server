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
	ingredient_mappings := []schemas.Ingredient_Mapping_With_Name{}
	rows, err := query_ingredient_mappings(db, reqData.Ingredient_ID)
	if err != nil && err == sql.ErrNoRows {
		log.Println("Get_Ingredient_Details | ingredient does not exist: ", err.Error())
		return utilities.Send_Error(c, "Ingredient does not exist", fiber.StatusBadRequest)
	}
	if err != nil {
		log.Println("error in querying query_ingredient_mappings: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		err = scan_ingredient(rows, &ingredient_mappings)
		if err != nil && err == sql.ErrNoRows {
			log.Println("Get_Ingredient_Details | error in scanning ingredient mapping: ", err.Error())
			return utilities.Send_Error(c, "Ingredient does not exist", fiber.StatusInternalServerError)
		}
		if err != nil {
			log.Println("Get_Ingredient_Details | error in scanning ingredient mapping: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
	}
	response.Ingredient_Mappings = ingredient_mappings
	temp_res := schemas.Res_Get_Ingredient_Mapping_Details{}
	if len(response.Ingredient_Mappings) < 1 {
		log.Println("Get_Ingredient_Details | ingredient does not exist")
		return utilities.Send_Error(c, "Ingredient does not exist", fiber.StatusBadRequest)
	}
	if len(response.Ingredient_Mappings) > 0 {
		row := query_ingredient_mapping(db, response.Ingredient_Mappings[0].ID)
		err = scan_ingredient_mapping(row, &temp_res)
		if err != nil && err == sql.ErrNoRows {
			log.Println("Get_Ingredient_Mapping_Details | error in scanning ingredient: ", err.Error())
			return utilities.Send_Error(c, "Ingredient does not exist", fiber.StatusBadRequest)
		}
		// querying ingredient images
		images, err := query_and_scan_food_images(db, response.Ingredient_Mappings[0].ID)
		if err != nil {
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		response.Ingredient_Images = images
	}
	response.Ingredient = temp_res.Ingredient
	response.Ingredient_Variant = temp_res.Ingredient_Variant
	response.Ingredient_Subvariant = temp_res.Ingredient_Subvariant
	response.Nutrient = temp_res.Nutrient
	response.Ingredient_Images = temp_res.Ingredient_Images
	return c.Status(fiber.StatusOK).JSON(response)
}

func query_ingredient_mappings(db *sql.DB, ingredient_id uint) (*sql.Rows, error) {
	rows, err := db.Query(`
		SELECT 
			ingredient_mapping.id,
			ingredient_mapping.ingredient_id,
			ingredient_mapping.ingredient_variant_id,
			ingredient_mapping.ingredient_subvariant_id,
			ingredient_mapping.nutrient_id,
			ingredient_variant.name,
			ingredient_subvariant.name
		FROM ingredient_mapping
		JOIN ingredient_variant ON ingredient_mapping.ingredient_variant_id = ingredient_variant.id
		JOIN ingredient_subvariant ON ingredient_mapping.ingredient_subvariant_id = ingredient_subvariant.id
		WHERE ingredient_mapping.ingredient_id = $1
		ORDER BY ingredient_variant.name ASC, ingredient_subvariant.name ASC`,
		ingredient_id,
	)
	return rows, err
}
func scan_ingredient(row *sql.Rows, ingredient_mappings *[]schemas.Ingredient_Mapping_With_Name) error {
	var ingredient_mapping = schemas.Ingredient_Mapping_With_Name{}
	if err := row.
		Scan(
			&ingredient_mapping.ID,
			&ingredient_mapping.Ingredient_Id,
			&ingredient_mapping.Ingredient_Variant_Id,
			&ingredient_mapping.Ingredient_Subvariant_Id,
			&ingredient_mapping.Nutrient_Id,
			&ingredient_mapping.Ingredient_Variant_Name,
			&ingredient_mapping.Ingredient_Subvariant_Name,
		); err != nil {
		return err
	}
	*ingredient_mappings = append(*ingredient_mappings, ingredient_mapping)
	return nil
}
