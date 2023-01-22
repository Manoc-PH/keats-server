package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/tracker"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Put_Intake(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Put_Intake | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Put_Intake)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Put_Intake | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	if reqData.Food_Id != 0 && reqData.Recipe_Id != 0 {
		log.Println("Put_Intake | Error: user sending recipe id and food id")
		return utilities.Send_Error(c, "only one food item id required, received 2", fiber.StatusBadRequest)
	}

	// Final response
	response_data := schemas.Res_Patch_Intake{}

	//* data processing
	if reqData.Food_Id != 0 {
		intake := models.Intake{}
		food := models.Food{}
		food_nutrient := models.Food_Nutrient{}
		macros_curr := models.Macros{}
		// TODO OPTIMIZATION: USE GO ROUTINES
		row := query_intake_food(reqData.Intake_ID, db)
		err = scan_intake_food(row, &intake, &food, &food_nutrient)
		if err == sql.ErrNoRows {
			return utilities.Send_Error(c, "intake not found", fiber.StatusBadRequest)
		}
		if err != nil {
			log.Println("Put_Intake | Error on scanning food: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		row = query_macros(db, owner_id)
		err = scan_macros(row, &macros_curr)
		if err != nil {
			log.Println("Put_Intake | Error on scanning macros: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		new_coins, new_xp, new_deductions := 0, 0, 0
		old_intake_macros := models.Macros{ID: macros_curr.ID, Account_Id: owner_id}
		calc_macros(&old_intake_macros, &food_nutrient, intake.Amount)
		new_intake_macros := models.Macros{ID: macros_curr.ID, Account_Id: owner_id}
		calc_macros(&new_intake_macros, &food_nutrient, reqData.Amount)
		// ! STILL UNSURE OF THIS CODE BLOCK'S STABILITY (like my emotions)
		if old_intake_macros.Calories != new_intake_macros.Calories {
			old_coins, old_xp, old_deductions := utilities.Calc_CnXP_On_Intake(float32(old_intake_macros.Calories), float32(macros_curr.Calories-old_intake_macros.Calories), float32(macros_curr.Max_Calories))
			new_coins, new_xp, new_deductions = utilities.Calc_CnXP_On_Intake(float32(new_intake_macros.Calories), float32(macros_curr.Calories-old_intake_macros.Calories), float32(macros_curr.Max_Calories))
			new_deductions = (new_deductions * -1) + old_deductions
			new_coins = (new_coins - old_coins) + new_deductions
			new_xp = (new_xp - old_xp) + new_deductions
		}

		macros_to_add := models.Macros{ID: macros_curr.ID, Account_Id: owner_id}
		calc_macros_update(&old_intake_macros, &new_intake_macros, &macros_to_add)

		new_intake := models.Intake{}
		new_intake = intake
		new_intake.Amount = reqData.Amount
		new_intake.Amount_Unit = reqData.Amount_Unit
		new_intake.Amount_Unit_Desc = reqData.Amount_Unit_Desc
		new_intake.Serving_Size = reqData.Serving_Size
		err = update_intake_macro_and_gamestat(db, &macros_to_add, new_coins, new_xp, &new_intake)
		if err != nil {
			log.Println("Put_Intake | Error on update_intake_macro_and_gamestat: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		response_data.Intake = new_intake
		response_data.Added_Coins_And_XP = schemas.Added_Coins_And_XP{Coins: new_coins, XP: new_xp}
		response_data.Added_Macros = schemas.Added_Macros{
			Calories: macros_to_add.Calories,
			Protein:  macros_to_add.Protein,
			Carbs:    macros_to_add.Carbs,
			Fats:     macros_to_add.Fats,
		}
		response_data.Food = food
	}
	// TODO ADD SUPPORT FOR RECIPES
	if reqData.Recipe_Id != 0 {
		return utilities.Send_Error(c, "recipes not yet supported", fiber.StatusBadRequest)
	}

	return c.Status(fiber.StatusOK).JSON(response_data)
}

func query_intake_food(intake_id uint, db *sql.DB) *sql.Row {
	row := db.QueryRow(`SELECT
			intake.id,
			intake.food_id,
			intake.amount,
			intake.amount_unit,
			intake.amount_unit_desc,
			intake.serving_size,
			food.id, food.name, food.name_ph, food.name_brand,
			food_nutrient.id,
			food_nutrient.amount,
			food_nutrient.amount_unit,
			food_nutrient.amount_unit_desc,
			food_nutrient.serving_size,
			food_nutrient.calories,
			food_nutrient.protein,
			food_nutrient.carbs,
			food_nutrient.fats
		FROM intake
		JOIN food ON intake.food_id = food.id
		JOIN food_nutrient ON food.food_nutrient_id = food_nutrient.id
		WHERE intake.id = $1`,
		intake_id,
	)
	return row
}
func scan_intake_food(row *sql.Row, intake *models.Intake, food *models.Food, food_nutrient *models.Food_Nutrient) error {
	if err := row.
		Scan(
			&intake.ID,
			&intake.Food_Id,
			&intake.Amount,
			&intake.Amount_Unit,
			&intake.Amount_Unit_Desc,
			&intake.Serving_Size,

			&food.ID,
			&food.Name,
			&food.Name_Ph,
			&food.Name_Brand,

			&food_nutrient.ID,
			&food_nutrient.Amount,
			&food_nutrient.Amount_Unit,
			&food_nutrient.Amount_Unit_Desc,
			&food_nutrient.Serving_Size,
			&food_nutrient.Calories,
			&food_nutrient.Protein,
			&food_nutrient.Carbs,
			&food_nutrient.Fats,
		); err != nil {
		return err
	}
	return nil
}
func calc_macros_update(old_macros *models.Macros, new_macros *models.Macros, macros_to_add *models.Macros) {
	macros_to_add.Calories = new_macros.Calories - old_macros.Calories
	macros_to_add.Protein = new_macros.Protein - old_macros.Protein
	macros_to_add.Carbs = new_macros.Carbs - old_macros.Carbs
	macros_to_add.Fats = new_macros.Fats - old_macros.Fats
}
func update_intake_macro_and_gamestat(db *sql.DB, macros_to_add *models.Macros, coins int, xp int, intake *models.Intake) error {
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = txn.Exec(
		`UPDATE macros SET
			calories = calories + $1,
			protein = protein + $2,
			carbs = carbs + $3,
			fats = fats + $4
		WHERE id = $5`,
		macros_to_add.Calories,
		macros_to_add.Protein,
		macros_to_add.Carbs,
		macros_to_add.Fats,
		macros_to_add.ID,
	)
	if err != nil {
		log.Println("update_intake_macro_and_gamestat (update macros) | Error: ", err.Error())
		return err
	}
	_, err = txn.Exec(
		`UPDATE account_game_stat SET coins = coins + $1, xp = xp + $2 WHERE account_id = $3`,
		coins, xp, macros_to_add.Account_Id,
	)
	if err != nil {
		log.Println("update_intake_macro_and_gamestat (update account_game_stat)| Error: ", err.Error())
		return err
	}
	_, err = txn.Exec(
		`UPDATE intake SET
			date_created = $1,
			amount = $2,
			amount_unit = $3,
			amount_unit_desc = $4,
			serving_size = $5
		WHERE id = $6`,
		time.Now(),
		intake.Amount,
		intake.Amount_Unit,
		intake.Amount_Unit_Desc,
		intake.Serving_Size,
		intake.ID,
	)
	if err != nil {
		log.Println("update_intake_macro_and_gamestat (update intake)| Error: ", err.Error())
		return err
	}
	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("update_intake_macro_and_gamestat (commit) | Error: ", err.Error())
		return err
	}
	return nil
}
