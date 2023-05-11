package handlers

import (
	"database/sql"
	"log"
	constants "server/constants/formats"
	"server/middlewares"
	schemas "server/schemas/tracker"
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
			COALESCE(intake.food_id, 0) as food_id,
			COALESCE(intake.recipe_id, 0) as recipe_id,
			intake.amount,
			intake.amount_unit,
			intake.amount_unit_desc,
			intake.serving_size,
			COALESCE(food.name, '') as food_name,
			COALESCE(food.name_ph, '') as food_name_ph,
			COALESCE(food.name_brand, '') as food_name_brand,
			COALESCE(food.food_nutrient_id, 0) as food_nutrient_id,
			COALESCE(food_nutrient.calories, 0) as food_nutrient_calories,
			COALESCE(food_nutrient.amount, 0) as food_nutrient_amount,
			COALESCE(food_nutrient.amount_unit, '') as food_nutrient_amount_unit,
			COALESCE(recipe.name, '') as recipe_name,
			COALESCE(recipe.name_owner, '') as recipe_name_owner
		FROM intake
		LEFT JOIN food ON intake.food_id = food.id
		LEFT JOIN food_nutrient ON food.food_nutrient_id = food_nutrient.id
		LEFT JOIN recipe ON intake.recipe_id = recipe.id
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
				&new_intake.Food_Id,
				&new_intake.Recipe_Id,

				&new_intake.Amount,
				&new_intake.Amount_Unit,
				&new_intake.Amount_Unit_Desc,
				&new_intake.Serving_Size,

				&new_intake.Food_Name,
				&new_intake.Food_Name_Ph,
				&new_intake.Food_Name_Brand,
				&new_intake.Food_Nutrient_Id,
				&new_intake.Food_Nutrient_Calories,
				&new_intake.Food_Nutrient_Amount,
				&new_intake.Food_Nutrient_Amount_Unit,
				&new_intake.Recipe_Name,
				&new_intake.Recipe_Name_Owner,
			); err != nil {
			log.Println("Get_Intakes | error in scanning intakes: ", err.Error())
			return nil, err
		}
		new_intake.Calories = (new_intake.Food_Nutrient_Calories / new_intake.Food_Nutrient_Amount) * new_intake.Amount
		intakes = append(intakes, new_intake)
	}
	return intakes, nil
}
