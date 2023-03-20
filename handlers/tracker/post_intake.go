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

func Post_Intake(c *fiber.Ctx, db *sql.DB) error {
	// TODO Insert new daily nutrient
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Intake | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Post_Intake)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Intake | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	if reqData.Food_Id != 0 && reqData.Recipe_Id != 0 {
		log.Println("Post_Intake | Error: user sending recipe id and food id")
		return utilities.Send_Error(c, "only one food item id required, received 2", fiber.StatusBadRequest)
	}

	// Final response
	response_data := schemas.Res_Post_Intake{}

	//* data processing
	if reqData.Food_Id != 0 {
		food := models.Food{}
		food_nutrient := models.Food_Nutrient{}
		d_nutrients_curr := models.Daily_Nutrients{}
		// TODO OPTIMIZATION: USE GO ROUTINES
		row := query_food(reqData.Food_Id, db)
		err = scan_food(row, &food, &food_nutrient)
		if err != nil {
			log.Println("Post_Intake | Error on scanning food: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		row = query_d_nutrients(db, owner_id)
		err = scan_d_nutrients(row, &d_nutrients_curr)
		if err != nil {
			log.Println("Post_Intake | Error on scanning d_nutrients: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		//! d_nutrients to add doesnt return the total d_nutrients
		d_nutrients_to_add := models.Daily_Nutrients{ID: d_nutrients_curr.ID, Account_Id: owner_id}
		new_intake := models.Intake{
			Account_Id:       owner_id,
			Date_Created:     time.Now(),
			Food_Id:          food.ID,
			Amount:           reqData.Amount,
			Amount_Unit:      reqData.Amount_Unit,
			Amount_Unit_Desc: reqData.Amount_Unit_Desc,
			Serving_Size:     reqData.Serving_Size,
		}
		calc_d_nutrients(&d_nutrients_to_add, &food_nutrient, reqData.Amount)
		coins, xp, deductions := utilities.Calc_CnXP_On_Intake(float32(d_nutrients_to_add.Calories), float32(d_nutrients_curr.Calories), float32(d_nutrients_curr.Max_Calories))
		coins = coins - deductions
		xp = xp - deductions
		err = save_intake_d_nutrients_and_gamestat(db, &d_nutrients_to_add, coins, xp, &new_intake)
		if err != nil {
			log.Println("Post_Intake | Error on save_intake_d_nutrients_and_gamestat: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		response_data.Intake = new_intake
		response_data.Added_Coins_And_XP = schemas.Added_Coins_And_XP{Coins: coins, XP: xp}
		response_data.Added_Daily_Nutrients = schemas.Added_Daily_Nutrients{
			Calories: d_nutrients_to_add.Calories,
			Protein:  d_nutrients_to_add.Protein,
			Carbs:    d_nutrients_to_add.Carbs,
			Fats:     d_nutrients_to_add.Fats,
		}
		response_data.Food = food
	}
	// TODO ADD SUPPORT FOR RECIPES
	if reqData.Recipe_Id != 0 {
		return utilities.Send_Error(c, "recipes not yet supported", fiber.StatusBadRequest)
	}
	return c.Status(fiber.StatusOK).JSON(response_data)
}

func query_food(food_id uint, db *sql.DB) *sql.Row {
	row := db.QueryRow(`SELECT
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
		FROM food
		JOIN food_nutrient ON food.food_nutrient_id = food_nutrient.id
		WHERE food.id = $1`,
		// casting timestamp to date
		food_id,
	)
	return row
}
func scan_food(row *sql.Row, food *models.Food, food_nutrient *models.Food_Nutrient) error {
	if err := row.
		Scan(
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
func calc_d_nutrients(d_nutrients_to_add *models.Daily_Nutrients, food_nutrient *models.Food_Nutrient, amount float32) {
	// TODO ADD HANDLER FOR DIFFERENT AMOUNT UNIT ||
	// TODO WRITE A CONVERTER THAT CHANGES THE food_nutrient AMOUNT VALUE TO GRAMS
	// if reqData.Amount_Unit != food_nutrient.Amount_Unit {}

	// Servings should be converted to amount in grams in the frontend
	amount_modifier := amount / food_nutrient.Amount
	d_nutrients_to_add.Calories = (food_nutrient.Calories * amount_modifier)
	d_nutrients_to_add.Protein = (food_nutrient.Protein * amount_modifier)
	d_nutrients_to_add.Carbs = (food_nutrient.Carbs * amount_modifier)
	d_nutrients_to_add.Fats = (food_nutrient.Fats * amount_modifier)
}
func save_intake_d_nutrients_and_gamestat(db *sql.DB, d_nutrients_to_add *models.Daily_Nutrients, coins int, xp int, intake *models.Intake) error {
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = txn.Exec(
		`UPDATE daily_nutrients SET
			calories = calories + $1,
			protein = protein + $2,
			carbs = carbs + $3,
			fats = fats + $4
		WHERE id = $5`,
		d_nutrients_to_add.Calories,
		d_nutrients_to_add.Protein,
		d_nutrients_to_add.Carbs,
		d_nutrients_to_add.Fats,
		d_nutrients_to_add.ID,
	)
	if err != nil {
		log.Println("save_intake_d_nutrients_and_gamestat (update d_nutrients) | Error: ", err.Error())
		return err
	}
	_, err = txn.Exec(
		`UPDATE account_game_stat SET coins = coins + $1, xp = xp + $2 WHERE account_id = $3`,
		coins, xp, d_nutrients_to_add.Account_Id,
	)
	if err != nil {
		log.Println("save_intake_d_nutrients_and_gamestat (update account_game_stat)| Error: ", err.Error())
		return err
	}
	if intake.Food_Id != 0 && intake.Recipe_Id == 0 {
		_, err = txn.Exec(
			`INSERT INTO intake (account_id, date_created, food_id, amount,	amount_unit, amount_unit_desc, serving_size)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			intake.Account_Id,
			intake.Date_Created,
			intake.Food_Id,
			intake.Amount,
			intake.Amount_Unit,
			intake.Amount_Unit_Desc,
			intake.Serving_Size,
		)
		if err != nil {
			log.Println("save_intake_d_nutrients_and_gamestat (insert intake)| Error: ", err.Error())
			return err
		}
	}
	if intake.Food_Id == 0 && intake.Recipe_Id != 0 {
		row := txn.QueryRow(
			`INSERT INTO intake (account_id, date_created, recipe_id, amount,	amount_unit, amount_unit_desc, serving_size)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			intake.Account_Id,
			intake.Date_Created,
			intake.Recipe_Id,
			intake.Amount,
			intake.Amount_Unit,
			intake.Amount_Unit_Desc,
			intake.Serving_Size,
		)
		err := row.Scan(&intake.ID)
		if err != nil {
			log.Println("save_intake_d_nutrients_and_gamestat (insert intake)| Error: ", err.Error())
			return err
		}
	}
	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("save_intake_d_nutrients_and_gamestat (commit) | Error: ", err.Error())
		return err
	}
	return nil
}
