package handlers

import (
	"database/sql"
	"log"
	constants "server/constants/formats"
	"server/middlewares"
	schemas "server/schemas/consumer/tracker"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Gets the intakes for the day
func Get_Intakes(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, Owner_Id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Intakes | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	// querying intakes
	intakes, err := query_and_scan_intakes(db, Owner_Id)
	// Intakes dont exist yet
	if err != nil && err == sql.ErrNoRows {
		log.Println("Get_Intakes | error in query_and_scan_intakes: ", err.Error())
		return utilities.Send_Error(c, "No intake found", fiber.StatusBadRequest)
	}
	// Server Error
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get_Intakes | error in scanning intakes: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(intakes)
}

func query_and_scan_intakes(db *sql.DB, user_id uuid.UUID) ([]schemas.Res_Get_Intakes, error) {
	rows, err := db.Query(`SELECT
			intake.id,
			intake.account_id,
			intake.date_created,
			COALESCE(intake.ingredient_mapping_id, 0) as ingredient_mapping_id,
			COALESCE(intake.food_id, 0) as food_id,
			intake.amount,
			intake.amount_unit,
			intake.amount_unit_desc,
			intake.serving_size,
			COALESCE(food.name, '') as food_name,
			COALESCE(food.name_ph, '') as food_name_ph,
			COALESCE(food.name_owner, '') as food_name_owner,
			COALESCE(ingredient.name, '') as ingredient_name,
			COALESCE(ingredient.name_ph, '') as ingredient_name_ph,
			COALESCE(ingredient_variant.name, '') as ingredient_variant_name,
			COALESCE(ingredient_variant.name_ph, '') as ingredient_variant_name_ph,
			COALESCE(ingredient_subvariant.name, '') as ingredient_subvariant_name,
			COALESCE(ingredient_subvariant.name_ph, '') as ingredient_subvariant_name_ph,
			COALESCE(ingredient.name_owner, '') as ingredient_name_owner
		FROM intake
		LEFT JOIN food ON intake.food_id = food.id
		LEFT JOIN ingredient_mapping ON intake.ingredient_mapping_id = ingredient_mapping.id
		LEFT JOIN ingredient ON ingredient_mapping.ingredient_id = ingredient.id
		LEFT JOIN ingredient_variant ON ingredient_mapping.ingredient_variant_id = ingredient_variant.id
		LEFT JOIN ingredient_subvariant ON ingredient_mapping.ingredient_subvariant_id = ingredient_subvariant.id
		WHERE intake.account_id = $1 AND intake.date_created >= $2
		ORDER BY intake.date_created DESC`,
		user_id, time.Now().Format(constants.YYYY_MM_DD),
	)
	if err != nil {
		log.Println("Get_Intakes | error in querying intakes: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	intakes := make([]schemas.Res_Get_Intakes, 0, 100)
	for rows.Next() {
		var new_intake = schemas.Res_Get_Intakes{}
		if err := rows.
			Scan(
				&new_intake.ID,
				&new_intake.Account_Id,
				&new_intake.Date_Created,
				&new_intake.Ingredient_Mapping_Id,
				&new_intake.Food_Id,

				&new_intake.Amount,
				&new_intake.Amount_Unit,
				&new_intake.Amount_Unit_Desc,
				&new_intake.Serving_Size,

				&new_intake.Food_Name,
				&new_intake.Food_Name_Ph,
				&new_intake.Food_Name_Owner,

				&new_intake.Ingredient_Name,
				&new_intake.Ingredient_Name_Ph,
				&new_intake.Ingredient_Variant_Name,
				&new_intake.Ingredient_Variant_Name_Ph,
				&new_intake.Ingredient_Subvariant_Name,
				&new_intake.Ingredient_Subvariant_Name_Ph,
				&new_intake.Ingredient_Name_Owner,
			); err != nil {
			log.Println("Get_Intakes | error in scanning intakes: ", err.Error())
			return nil, err
		}
		// new_intake.Calories = (new_intake.Food_Nutrient_Calories / new_intake.Food_Nutrient_Amount) * new_intake.Amount
		intakes = append(intakes, new_intake)
	}
	return intakes, nil
}
