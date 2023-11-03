package handlers

import (
	"database/sql"
	"log"
	constants "server/constants/formats"
	"server/middlewares"
	schemas "server/schemas/consumer/tracker"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Gets the summary of daily nutrients through a date range
func Get_Common_Intakes(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Common_Intakes | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Get_Common_Intakes)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Common_Intakes | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// querying new_intake
	common_intakes, err := query_and_scan_common_intakes(db, owner_id, reqData)
	if err != nil && err != sql.ErrNoRows {
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(*common_intakes)
}

func query_and_scan_common_intakes(db *sql.DB, owner_id uuid.UUID, reqData *schemas.Req_Get_Common_Intakes) (*schemas.Res_Get_Common_Intakes, error) {
	rows, err := db.Query(`
			SELECT
				intake.food_id,
				intake.ingredient_mapping_id,
				COUNT(intake.ingredient_mapping_id) AS ingredient_count,
				COUNT(intake.food_id) AS food_count,
				ingredient.id, coalesce(ingredient.name, ''), coalesce(ingredient.name_ph, ''), coalesce(ingredient.name_owner, ''),
				ingredient_variant.id, coalesce(ingredient_variant.name, ''), coalesce(ingredient_variant.name_ph, ''), 
				ingredient_subvariant.id, coalesce(ingredient_subvariant.name, ''), coalesce(ingredient_subvariant.name_ph, ''),
				coalesce(food.name, ''), coalesce(food.name_ph, ''), coalesce(food.name_owner, '')
			FROM intake
			LEFT JOIN ingredient_mapping ON intake.ingredient_mapping_id = ingredient_mapping.id   
			LEFT JOIN ingredient ON ingredient_mapping.ingredient_id = ingredient.id
			LEFT JOIN ingredient_variant ON ingredient_mapping.ingredient_variant_id = ingredient_variant.id
			LEFT JOIN ingredient_subvariant ON ingredient_mapping.ingredient_subvariant_id = ingredient_subvariant.id
			LEFT JOIN food ON intake.food_id = food.id
			WHERE intake.account_id = $1
			AND intake.date_created BETWEEN $2 AND $3
			GROUP BY
				intake.food_id,
				intake.ingredient_mapping_id,
				ingredient.id,
				ingredient_variant.id,
				ingredient_subvariant.id,
				food.id HAVING COUNT(*) >= 1
			ORDER BY COUNT(*) DESC LIMIT 15`,
		owner_id, reqData.Start_Date.Format(constants.YYYY_MM_DD), reqData.End_Date.Format(constants.YYYY_MM_DD),
	)
	if err != nil {
		log.Println("Get_Common_Intakes | error in querying common intakes: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	response := schemas.Res_Get_Common_Intakes{}
	for rows.Next() {
		var new_intake = schemas.Intake_Details{}
		if err := rows.
			Scan(
				&new_intake.Food_ID,
				&new_intake.Ingredient_Mapping_ID,
				&new_intake.Ingredient_Count,
				&new_intake.Food_Count,
				&new_intake.Ingredient_ID,
				&new_intake.Ingredient_Name,
				&new_intake.Ingredient_Name_Ph,
				&new_intake.Ingredient_Name_Owner,
				&new_intake.Ingredient_Variant_ID,
				&new_intake.Ingredient_Variant_Name,
				&new_intake.Ingredient_Variant_Name_Ph,
				&new_intake.Ingredient_Subvariant_ID,
				&new_intake.Ingredient_Subvariant_Name,
				&new_intake.Ingredient_Subvariant_Name_Ph,
				&new_intake.Food_Name,
				&new_intake.Food_Name_Ph,
				&new_intake.Food_Name_Owner,
			); err != nil {
			log.Println("Get_Common_Intakes | error in scanning new_intake: ", err.Error())
			return nil, err
		}
		response.Intakes = append(response.Intakes, new_intake)
	}
	return &response, nil
}
