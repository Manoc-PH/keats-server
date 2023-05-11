package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/tracker"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TODO Update this handler, make sure to only query either food or recipe and not join them into a single struct

// Gets the details of the intake
func Get_Intake_Details(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, Owner_Id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Intake_Details | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Intake_Details)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Intake_Details | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	response := schemas.Res_Get_Intake_Details{}
	// querying intake
	row := query_intake(db, Owner_Id, reqData.Intake_ID)
	// scanning intake
	err = scan_intake(row, &response)
	if err != nil && err == sql.ErrNoRows {
		log.Println("Get_Intake_Details | error in scanning intake: ", err.Error())
		return utilities.Send_Error(c, "Intake does not exist", fiber.StatusBadRequest)
	}
	// Server Error
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get_Intake_Details | error in scanning intake: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func query_intake(db *sql.DB, user_id uuid.UUID, intake_id uint) *sql.Row {
	row := db.QueryRow(`SELECT
			intake.id,
			intake.account_id,
			intake.date_created,
			COALESCE(intake.food_id, 0) as food_id,
			COALESCE(intake.recipe_id, 0) as recipe_id,
			intake.amount,
			intake.amount_unit,
			intake.amount_unit_desc,
			intake.serving_size,
			COALESCE(food.name, '') as name,
			COALESCE(food.name_ph, '') as name_ph,
			COALESCE(food.name_brand, '') as name_brand,
			COALESCE(food.food_nutrient_id, 0) as food_nutrient_id,
			--	FOOD NUTRIENT
			COALESCE(food_nutrient.amount, 0) as amount,
			COALESCE(food_nutrient.amount_unit, '') as amount_unit,
			COALESCE(food_nutrient.amount_unit_desc, '') as amount_unit_desc,
			COALESCE(food_nutrient.serving_size, 0) as serving_size,
			COALESCE(food_nutrient.calories, 0) as calories,
			COALESCE(food_nutrient.protein, 0) as protein,
			COALESCE(food_nutrient.carbs, 0) as carbs,
			COALESCE(food_nutrient.fats, 0) as fats,
			COALESCE(food_nutrient.trans_fat, 0) as trans_fat,
			COALESCE(food_nutrient.saturated_fat, 0) as saturated_fat,
			COALESCE(food_nutrient.sugars, 0) as sugars,
			COALESCE(food_nutrient.sodium, 0) as sodium,
			--	RECIPE
			COALESCE(recipe.name, '') as name,
			COALESCE(recipe.name_owner, '') as name_owner
		FROM intake
		LEFT JOIN food ON intake.food_id = food.id
		LEFT JOIN food_nutrient ON food.food_nutrient_id = food_nutrient.id
		LEFT JOIN recipe ON intake.recipe_id = recipe.id
		WHERE intake.account_id = $1 AND intake.id = $2`,
		user_id, intake_id,
	)
	return row
}

func scan_intake(row *sql.Row, intake *schemas.Res_Get_Intake_Details) error {
	err := row.Scan(
		&intake.ID,
		&intake.Account_Id,
		&intake.Date_Created,
		&intake.Food_Id,
		&intake.Recipe_Id,

		&intake.Amount,
		&intake.Amount_Unit,
		&intake.Amount_Unit_Desc,
		&intake.Serving_Size,

		&intake.Food.Details.Name,
		&intake.Food.Details.Name_Ph,
		&intake.Food.Details.Name_Brand,
		&intake.Food.Details.Food_Nutrient_Id,
		&intake.Food.Details.Amount,
		&intake.Food.Details.Amount_Unit,
		&intake.Food.Details.Amount_Unit_Desc,
		&intake.Food.Details.Serving_Size,
		&intake.Food.Details.Calories,
		&intake.Food.Details.Protein,
		&intake.Food.Details.Carbs,
		&intake.Food.Details.Fats,
		&intake.Food.Details.Trans_Fat,
		&intake.Food.Details.Saturated_Fat,
		&intake.Food.Details.Sugars,
		&intake.Food.Details.Sodium,

		&intake.Recipe.Name,
		&intake.Recipe.Name_Owner,
	)
	return err
}
