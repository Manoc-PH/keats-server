package handlers

import (
	"database/sql"
	"log"
	"server/constants"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/tracker"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	if reqData.Food_Id != constants.Empty_UUID && reqData.Ingredient_Mapping_Id != constants.Empty_UUID {
		log.Println("Put_Intake | Error: user sending recipe id and food id")
		return utilities.Send_Error(c, "only one food item id required, received 2", fiber.StatusBadRequest)
	}

	// Final response
	response_data := schemas.Res_Put_Intake{}

	//* data processing
	if reqData.Ingredient_Mapping_Id != constants.Empty_UUID {
		intake := models.Intake{}
		new_nutrient := models.Nutrient{}
		old_nutrient := models.Nutrient{}
		daily_nutrients := models.Daily_Nutrients{Account_Id: owner_id}
		// TODO OPTIMIZATION: USE GO ROUTINES
		// Querrying intake
		row := query_intake(db, owner_id, reqData.Intake_ID)
		err = scan_intake(row, &intake)
		if err == sql.ErrNoRows {
			return utilities.Send_Error(c, "intake not found", fiber.StatusBadRequest)
		}
		if err != nil {
			log.Println("Put_Intake | Error on scanning food: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		// Not allowing user to edit intake from past
		is_intake_today := check_if_date_is_today(intake.Date_Created, time.Now())
		if !is_intake_today {
			log.Println("Put_Intake | Error: User trying to edit old intake")
			return utilities.Send_Error(c, "cannot edit intake from more than a day ago", fiber.StatusBadRequest)
		}
		// Querying old nutrients of ingredient
		row = query_ingredient_nutrient(intake.Ingredient_Mapping_Id, db)
		err = scan_nutrient(row, &old_nutrient)
		if err != nil {
			log.Println("Put_Intake | Error on scanning nutrients: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		// TODO Optimize this area, I dont think it needs to be called twice
		// Querying new nutrients of ingredient
		row = query_ingredient_nutrient(reqData.Ingredient_Mapping_Id, db)
		err = scan_nutrient(row, &new_nutrient)
		if err != nil {
			log.Println("Put_Intake | Error on scanning nutrients: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		// Getting daily nutrients
		row = query_daily_nutrients(db, owner_id)
		err = scan_daily_nutrients(row, &daily_nutrients)
		if err != nil {
			log.Println("Put_Intake | Error on scanning daily_nutrients: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		old_intake_d_nutrients := models.Nutrient{}
		calc_nutrients(&old_intake_d_nutrients, &old_nutrient, intake.Amount)
		new_intake_d_nutrients := models.Nutrient{}
		calc_nutrients(&new_intake_d_nutrients, &new_nutrient, reqData.Amount)
		daily_nutrients_to_add := models.Daily_Nutrients{ID: daily_nutrients.ID}
		calc_daily_nutrients_update(&old_intake_d_nutrients, &new_intake_d_nutrients, &daily_nutrients_to_add)

		new_intake := intake
		new_intake.Amount = reqData.Amount
		new_intake.Amount_Unit = reqData.Amount_Unit
		new_intake.Amount_Unit_Desc = reqData.Amount_Unit_Desc
		new_intake.Serving_Size = reqData.Serving_Size
		new_intake.Ingredient_Mapping_Id = reqData.Ingredient_Mapping_Id
		txn, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		// Updating Intake
		err = update_intake(txn, &new_intake)
		if err != nil {
			log.Println("Put_Intake | Error on update_intake: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		// Updating Daily Nutrients
		err = update_daily_nutrients(txn, &daily_nutrients_to_add)
		if err != nil {
			log.Println("Put_Intake | Error on update_daily_nutrients: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		err = txn.Commit()
		if err != nil {
			txn.Rollback()
			log.Println("update_intake_d_nutrients_and_gamestat (commit) | Error: ", err.Error())
			return err
		}
		response_data.Added_Daily_Nutrients.Calories = daily_nutrients_to_add.Calories
		response_data.Added_Daily_Nutrients.Protein = daily_nutrients_to_add.Protein
		response_data.Added_Daily_Nutrients.Carbs = daily_nutrients_to_add.Carbs
		response_data.Added_Daily_Nutrients.Fats = daily_nutrients_to_add.Fats
		response_data.Added_Daily_Nutrients.Trans_Fat = daily_nutrients_to_add.Trans_Fat
		response_data.Added_Daily_Nutrients.Saturated_Fat = daily_nutrients_to_add.Saturated_Fat
		response_data.Added_Daily_Nutrients.Sugars = daily_nutrients_to_add.Sugars
		response_data.Added_Daily_Nutrients.Fiber = daily_nutrients_to_add.Fiber
		response_data.Added_Daily_Nutrients.Sodium = daily_nutrients_to_add.Sodium
		response_data.Added_Daily_Nutrients.Iron = daily_nutrients_to_add.Iron
		response_data.Added_Daily_Nutrients.Calcium = daily_nutrients_to_add.Calcium

		ingredient_mapping := &schemas.Ingredient_Mapping_Schema{}
		// Getting ingredient data
		row = query_ingredient(reqData.Ingredient_Mapping_Id, db)
		err = scan_ingredient(row, ingredient_mapping)
		if err != nil {
			log.Println("Post_Intake | Error on scanning ingredient: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}

		response_data.Ingredient = ingredient_mapping
		response_data.Intake = new_intake

	}
	if reqData.Food_Id != constants.Empty_UUID {
		intake := models.Intake{}
		new_nutrient := models.Nutrient{}
		old_nutrient := models.Nutrient{}
		daily_nutrients := models.Daily_Nutrients{Account_Id: owner_id}
		// Querrying intake
		row := query_intake(db, owner_id, reqData.Intake_ID)
		err = scan_intake(row, &intake)
		if err == sql.ErrNoRows {
			return utilities.Send_Error(c, "intake not found", fiber.StatusBadRequest)
		}
		if err != nil {
			log.Println("Put_Intake | Error on scanning food: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		// Not allowing user to edit intake from past
		is_intake_today := check_if_date_is_today(intake.Date_Created, time.Now())
		if !is_intake_today {
			log.Println("Put_Intake | Error: User trying to edit old intake")
			return utilities.Send_Error(c, "cannot edit intake from more than a day ago", fiber.StatusBadRequest)
		}
		// TODO Optimize this area, I dont think it needs to be called twice
		// Querying old nutrients of ingredient
		row = query_food_nutrient(intake.Food_Id, db)
		err = scan_nutrient(row, &old_nutrient)
		if err != nil {
			log.Println("Put_Intake | Error on scanning nutrients: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		// Querying new nutrients of ingredient
		row = query_food_nutrient(reqData.Food_Id, db)
		err = scan_nutrient(row, &new_nutrient)
		if err != nil {
			log.Println("Put_Intake | Error on scanning nutrients: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		// Getting daily nutrients
		row = query_daily_nutrients(db, owner_id)
		err = scan_daily_nutrients(row, &daily_nutrients)
		if err != nil {
			log.Println("Put_Intake | Error on scanning daily_nutrients: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		old_intake_d_nutrients := models.Nutrient{}
		calc_nutrients(&old_intake_d_nutrients, &old_nutrient, intake.Amount)
		new_intake_d_nutrients := models.Nutrient{}
		calc_nutrients(&new_intake_d_nutrients, &new_nutrient, reqData.Amount)
		daily_nutrients_to_add := models.Daily_Nutrients{ID: daily_nutrients.ID}
		calc_daily_nutrients_update(&old_intake_d_nutrients, &new_intake_d_nutrients, &daily_nutrients_to_add)

		new_intake := intake
		new_intake.Amount = reqData.Amount
		new_intake.Amount_Unit = reqData.Amount_Unit
		new_intake.Amount_Unit_Desc = reqData.Amount_Unit_Desc
		new_intake.Serving_Size = reqData.Serving_Size
		new_intake.Food_Id = reqData.Food_Id
		txn, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		// Updating Intake
		err = update_intake(txn, &new_intake)
		if err != nil {
			log.Println("Put_Intake | Error on update_intake: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		// Updating Daily Nutrients
		err = update_daily_nutrients(txn, &daily_nutrients_to_add)
		if err != nil {
			log.Println("Put_Intake | Error on update_daily_nutrients: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		err = txn.Commit()
		if err != nil {
			txn.Rollback()
			log.Println("update_intake_d_nutrients_and_gamestat (commit) | Error: ", err.Error())
			return err
		}
		response_data.Added_Daily_Nutrients.Calories = daily_nutrients_to_add.Calories
		response_data.Added_Daily_Nutrients.Protein = daily_nutrients_to_add.Protein
		response_data.Added_Daily_Nutrients.Carbs = daily_nutrients_to_add.Carbs
		response_data.Added_Daily_Nutrients.Fats = daily_nutrients_to_add.Fats
		response_data.Added_Daily_Nutrients.Trans_Fat = daily_nutrients_to_add.Trans_Fat
		response_data.Added_Daily_Nutrients.Saturated_Fat = daily_nutrients_to_add.Saturated_Fat
		response_data.Added_Daily_Nutrients.Sugars = daily_nutrients_to_add.Sugars
		response_data.Added_Daily_Nutrients.Fiber = daily_nutrients_to_add.Fiber
		response_data.Added_Daily_Nutrients.Sodium = daily_nutrients_to_add.Sodium
		response_data.Added_Daily_Nutrients.Iron = daily_nutrients_to_add.Iron
		response_data.Added_Daily_Nutrients.Calcium = daily_nutrients_to_add.Calcium

		food_mapping := &schemas.Food_Mapping_Schema{}
		// Getting food data
		row = query_food_and_nutrient(reqData.Food_Id, db)
		err = scan_food_and_nutrient(row, &food_mapping.Food, &food_mapping.Nutrient)
		if err != nil {
			log.Println("Post_Intake | Error on scanning food: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}

		response_data.Food = food_mapping
		response_data.Intake = new_intake
	}

	return c.Status(fiber.StatusOK).JSON(response_data)
}

func query_ingredient_nutrient(ingredient_mapping_id uuid.UUID, db *sql.DB) *sql.Row {
	row := db.QueryRow(`SELECT
			nutrient.id,
			nutrient.amount,
			coalesce(nutrient.amount_unit, ''),
			coalesce(nutrient.amount_unit_desc, ''),
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
		JOIN nutrient ON ingredient_mapping.nutrient_id = nutrient.id
		WHERE ingredient_mapping.id = $1`,
		// casting timestamp to date
		ingredient_mapping_id,
	)
	return row
}
func query_food_nutrient(food_id uuid.UUID, db *sql.DB) *sql.Row {
	row := db.QueryRow(`SELECT
			nutrient.id,
			nutrient.amount,
			coalesce(nutrient.amount_unit, ''),
			coalesce(nutrient.amount_unit_desc, ''),
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
		FROM food
		JOIN nutrient ON food.nutrient_id = nutrient.id
		WHERE food.id = $1`,
		// casting timestamp to date
		food_id,
	)
	return row
}
func scan_nutrient(row *sql.Row, nutrient *models.Nutrient) error {
	if err := row.
		Scan(
			&nutrient.ID,
			&nutrient.Amount,
			&nutrient.Amount_Unit,
			&nutrient.Amount_Unit_Desc,
			&nutrient.Serving_Size,
			&nutrient.Calories,
			&nutrient.Protein,
			&nutrient.Carbs,
			&nutrient.Fats,
			&nutrient.Trans_Fat,
			&nutrient.Saturated_Fat,
			&nutrient.Sugars,
			&nutrient.Fiber,
			&nutrient.Sodium,
			&nutrient.Iron,
			&nutrient.Calcium,
		); err != nil {
		return err
	}
	return nil
}
func calc_daily_nutrients_update(old_d_nutrients *models.Nutrient, new_d_nutrients *models.Nutrient, d_nutrients_to_add *models.Daily_Nutrients) {
	d_nutrients_to_add.Calories = new_d_nutrients.Calories - old_d_nutrients.Calories
	d_nutrients_to_add.Protein = new_d_nutrients.Protein - old_d_nutrients.Protein
	d_nutrients_to_add.Carbs = new_d_nutrients.Carbs - old_d_nutrients.Carbs
	d_nutrients_to_add.Fats = new_d_nutrients.Fats - old_d_nutrients.Fats
	d_nutrients_to_add.Trans_Fat = new_d_nutrients.Trans_Fat - old_d_nutrients.Trans_Fat
	d_nutrients_to_add.Saturated_Fat = new_d_nutrients.Saturated_Fat - old_d_nutrients.Saturated_Fat
	d_nutrients_to_add.Sugars = new_d_nutrients.Sugars - old_d_nutrients.Sugars
	d_nutrients_to_add.Fiber = new_d_nutrients.Fiber - old_d_nutrients.Fiber
	d_nutrients_to_add.Sodium = new_d_nutrients.Sodium - old_d_nutrients.Sodium
	d_nutrients_to_add.Iron = new_d_nutrients.Iron - old_d_nutrients.Iron
	d_nutrients_to_add.Calcium = new_d_nutrients.Calcium - old_d_nutrients.Calcium
}
func update_intake(txn *sql.Tx, intake *models.Intake) error {
	if intake.Ingredient_Mapping_Id != constants.Empty_UUID {
		_, err := txn.Exec(
			`UPDATE intake SET 
				amount = $1,
				amount_unit = $2,
				amount_unit_desc = $3,
				serving_size = $4,
				ingredient_mapping_id = $5
			WHERE id = $6`,
			intake.Amount,
			intake.Amount_Unit,
			intake.Amount_Unit_Desc,
			intake.Serving_Size,
			intake.Ingredient_Mapping_Id,
			intake.ID,
		)
		if err != nil {
			log.Println("update_intake | Error: ", err.Error())
			return err
		}
	}
	if intake.Food_Id != constants.Empty_UUID {
		_, err := txn.Exec(
			`UPDATE intake SET 
				amount = $1,
				amount_unit = $2,
				amount_unit_desc = $3,
				serving_size = $4, 
				food_id = $5
			WHERE id = $6`,
			intake.Amount,
			intake.Amount_Unit,
			intake.Amount_Unit_Desc,
			intake.Serving_Size,
			intake.Food_Id,
			intake.ID,
		)
		if err != nil {
			log.Println("update_intake | Error: ", err.Error())
			return err
		}
	}
	return nil
}
func update_daily_nutrients(txn *sql.Tx, daily_nutrients_to_add *models.Daily_Nutrients) error {
	_, err := txn.Exec(
		`UPDATE daily_nutrients SET
			calories = calories + $1,
			protein = protein + $2,
			carbs = carbs + $3,
			fats = fats + $4,
			trans_fat = trans_fat + $5,
			saturated_fat = saturated_fat + $6,
			sugars = sugars + $7,
			fiber = fiber + $8,
			sodium = sodium + $9,
			iron = iron + $10,
			calcium = calcium + $11
			WHERE id = $12`,
		daily_nutrients_to_add.Calories,
		daily_nutrients_to_add.Protein,
		daily_nutrients_to_add.Carbs,
		daily_nutrients_to_add.Fats,
		daily_nutrients_to_add.Trans_Fat,
		daily_nutrients_to_add.Saturated_Fat,
		daily_nutrients_to_add.Sugars,
		daily_nutrients_to_add.Fiber,
		daily_nutrients_to_add.Sodium,
		daily_nutrients_to_add.Iron,
		daily_nutrients_to_add.Calcium,
		daily_nutrients_to_add.ID,
	)
	if err != nil {
		log.Println("update_daily_nutrients | Error: ", err.Error())
		return err
	}
	return nil
}
func check_if_date_is_today(a time.Time, b time.Time) bool {
	if a.Day() == b.Day() && a.Month() == b.Month() && a.Year() == b.Year() {
		return true
	}
	return false
}

// func update_intake_d_nutrients_and_gamestat(db *sql.DB, d_nutrients_to_add *models.Nutrient, coins int, xp int, intake *models.Intake) error {
// 	txn, err := db.Begin()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	_, err = txn.Exec(
// 		`UPDATE daily_nutrients SET
// 			calories = calories + $1,
// 			protein = protein + $2,
// 			carbs = carbs + $3,
// 			fats = fats + $4
// 		WHERE id = $5`,
// 		d_nutrients_to_add.Calories,
// 		d_nutrients_to_add.Protein,
// 		d_nutrients_to_add.Carbs,
// 		d_nutrients_to_add.Fats,
// 		d_nutrients_to_add.ID,
// 	)
// 	if err != nil {
// 		log.Println("update_intake_d_nutrients_and_gamestat (update daily_nutrients) | Error: ", err.Error())
// 		return err
// 	}
// 	_, err = txn.Exec(
// 		`UPDATE account_game_stat SET coins = coins + $1, xp = xp + $2 WHERE account_id = $3`,
// 		coins, xp, intake.Account_Id,
// 	)
// 	if err != nil {
// 		log.Println("update_intake_d_nutrients_and_gamestat (update account_game_stat)| Error: ", err.Error())
// 		return err
// 	}
// 	_, err = txn.Exec(
// 		`UPDATE intake SET
// 			date_created = $1,
// 			amount = $2,
// 			amount_unit = $3,
// 			amount_unit_desc = $4,
// 			serving_size = $5
// 		WHERE id = $6`,
// 		time.Now(),
// 		intake.Amount,
// 		intake.Amount_Unit,
// 		intake.Amount_Unit_Desc,
// 		intake.Serving_Size,
// 		intake.ID,
// 	)
// 	if err != nil {
// 		log.Println("update_intake_d_nutrients_and_gamestat (update intake)| Error: ", err.Error())
// 		return err
// 	}
// 	err = txn.Commit()
// 	if err != nil {
// 		txn.Rollback()
// 		log.Println("update_intake_d_nutrients_and_gamestat (commit) | Error: ", err.Error())
// 		return err
// 	}
// 	return nil
// }
