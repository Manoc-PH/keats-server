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
	if reqData.Food_Id != 0 && reqData.Ingredient_Mapping_Id != 0 {
		log.Println("Post_Intake | Error: user sending recipe id and food id")
		return utilities.Send_Error(c, "only one food item id required, received 2", fiber.StatusBadRequest)
	}

	// Final response
	response_data := schemas.Res_Post_Intake{}

	//* data processing
	// *Ingredient
	if reqData.Ingredient_Mapping_Id != 0 {
		ingredient_mapping := schemas.Ingredient_Mapping_Schema{}
		daily_nutrients := models.Daily_Nutrients{Account_Id: owner_id}
		nutrients_to_add := models.Nutrient{}
		// Getting ingredient data
		row := query_ingredient(reqData.Ingredient_Mapping_Id, db)
		err = scan_ingredient(row, &ingredient_mapping)
		if err != nil {
			log.Println("Post_Intake | Error on scanning ingredient: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		// Getting daily nutrients
		row = query_daily_nutrients(db, owner_id)
		err = scan_daily_nutrients(row, &daily_nutrients)
		if err != nil {
			log.Println("Post_Intake | Error on scanning daily nutrients: ", err.Error())
			if err.Error() == sql.ErrNoRows.Error() {
				err = generate_daily_nutrients(db, owner_id, &daily_nutrients)
				if err != nil {
					log.Println("Post_Intake | error in generate_daily_nutrients: ", err.Error())
					return utilities.Send_Error(c, "An error occured in getting daily nutrients", fiber.StatusInternalServerError)
				}
			} else {
				return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
			}
		}
		// Calculating nutrients and saving to daily nutrients
		calc_nutrients(&nutrients_to_add, &ingredient_mapping.Nutrient, reqData.Amount)
		nutrients_to_add.ID = daily_nutrients.ID
		new_intake := models.Intake{
			Account_Id:            owner_id,
			Date_Created:          time.Now(),
			Ingredient_Mapping_Id: reqData.Ingredient_Mapping_Id,
			Amount:                reqData.Amount,
			Amount_Unit:           reqData.Amount_Unit,
			Amount_Unit_Desc:      reqData.Amount_Unit_Desc,
			Serving_Size:          reqData.Serving_Size,
		}
		// Saving acutal intake
		txn, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		err = save_daily_nutrients(txn, &nutrients_to_add)
		if err != nil {
			return utilities.Send_Error(c, "An error occured in saving daily nutrients", fiber.StatusInternalServerError)
		}
		err = save_intake_ingredient(txn, &new_intake)
		if err != nil {
			return utilities.Send_Error(c, "An error occured in saving intake", fiber.StatusInternalServerError)
		}
		err = txn.Commit()
		if err != nil {
			txn.Rollback()
			log.Println("save_intake_d_nutrients_and_gamestat (commit) | Error: ", err.Error())
			return err
		}
		response_data.Added_Daily_Nutrients = nutrients_to_add
		response_data.Ingredient = ingredient_mapping
		response_data.Intake = new_intake
	}

	// *Food
	// if reqData.Food_Id != 0 {
	// 	food := models.Food{}
	// 	food_nutrient := models.Food_Nutrient{}
	// 	d_nutrients_curr := models.Daily_Nutrients{}
	// 	// TODO OPTIMIZATION: USE GO ROUTINES
	// 	row := query_food(reqData.Food_Id, db)
	// 	err = scan_food(row, &food, &food_nutrient)
	// 	if err != nil {
	// 		log.Println("Post_Intake | Error on scanning food: ", err.Error())
	// 		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	// 	}
	// 	row = query_d_nutrients(db, owner_id)
	// 	err = scan_d_nutrients(row, &d_nutrients_curr)
	// 	if err != nil {
	// 		log.Println("Post_Intake | Error on scanning d_nutrients: ", err.Error())
	// 		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	// 	}
	// 	//! d_nutrients to add doesnt return the total d_nutrients
	// 	d_nutrients_to_add := models.Daily_Nutrients{ID: d_nutrients_curr.ID, Account_Id: owner_id}
	// 	new_intake := models.Intake{
	// 		Account_Id:       owner_id,
	// 		Date_Created:     time.Now(),
	// 		Food_Id:          food.ID,
	// 		Amount:           reqData.Amount,
	// 		Amount_Unit:      reqData.Amount_Unit,
	// 		Amount_Unit_Desc: reqData.Amount_Unit_Desc,
	// 		Serving_Size:     reqData.Serving_Size,
	// 	}
	// 	calc_d_nutrients(&d_nutrients_to_add, &food_nutrient, reqData.Amount)
	// 	coins, xp, deductions := utilities.Calc_CnXP_On_Intake(float32(d_nutrients_to_add.Calories), float32(d_nutrients_curr.Calories), float32(d_nutrients_curr.Max_Calories))
	// 	coins = coins - deductions
	// 	xp = xp - deductions
	// 	err = save_intake_d_nutrients_and_gamestat(db, &d_nutrients_to_add, coins, xp, &new_intake)
	// 	if err != nil {
	// 		log.Println("Post_Intake | Error on save_intake_d_nutrients_and_gamestat: ", err.Error())
	// 		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	// 	}
	// 	response_data.Intake = new_intake
	// 	// response_data.Added_Coins_And_XP = schemas.Added_Coins_And_XP{Coins: coins, XP: xp}
	// 	response_data.Added_Daily_Nutrients = schemas.Added_Daily_Nutrients{
	// 		Calories: d_nutrients_to_add.Calories,
	// 		Protein:  d_nutrients_to_add.Protein,
	// 		Carbs:    d_nutrients_to_add.Carbs,
	// 		Fats:     d_nutrients_to_add.Fats,
	// 	}
	// 	response_data.Food = food
	// }

	return c.Status(fiber.StatusOK).JSON(response_data)
}

func query_ingredient(ingredient_mapping_id uint, db *sql.DB) *sql.Row {
	row := db.QueryRow(`SELECT
			ingredient.id, ingredient.name, coalesce(ingredient.name_ph, ''), ingredient.name_brand,
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
func scan_ingredient(row *sql.Row, ingredient_mapping *schemas.Ingredient_Mapping_Schema) error {
	if err := row.
		Scan(
			&ingredient_mapping.Ingredient.ID,
			&ingredient_mapping.Ingredient.Name,
			&ingredient_mapping.Ingredient.Name_Ph,
			&ingredient_mapping.Ingredient.Name_Brand,

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
func calc_nutrients(nutrients_to_add *models.Nutrient, nutrient *models.Nutrient, amount float32) {
	// TODO ADD HANDLER FOR DIFFERENT AMOUNT UNIT ||
	// TODO WRITE A CONVERTER THAT CHANGES THE nutrient AMOUNT VALUE TO GRAMS
	// if reqData.Amount_Unit != nutrient.Amount_Unit {}

	// Servings should be converted to amount in grams in the frontend
	amount_modifier := amount / nutrient.Amount
	nutrients_to_add.Calories = (nutrient.Calories * amount_modifier)
	nutrients_to_add.Protein = (nutrient.Protein * amount_modifier)
	nutrients_to_add.Carbs = (nutrient.Carbs * amount_modifier)
	nutrients_to_add.Fats = (nutrient.Fats * amount_modifier)
	nutrients_to_add.Trans_Fat = (nutrient.Trans_Fat * amount_modifier)
	nutrients_to_add.Saturated_Fat = (nutrient.Saturated_Fat * amount_modifier)
	nutrients_to_add.Sugars = (nutrient.Sugars * amount_modifier)
	nutrients_to_add.Fiber = (nutrient.Fiber * amount_modifier)
	nutrients_to_add.Sodium = (nutrient.Sodium * amount_modifier)
	nutrients_to_add.Iron = (nutrient.Iron * amount_modifier)
	nutrients_to_add.Calcium = (nutrient.Calcium * amount_modifier)
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

			&food_nutrient.ID,
		); err != nil {
		return err
	}
	return nil
}
func save_daily_nutrients(txn *sql.Tx, d_nutrients_to_add *models.Nutrient) error {
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
		d_nutrients_to_add.Calories,
		d_nutrients_to_add.Protein,
		d_nutrients_to_add.Carbs,
		d_nutrients_to_add.Fats,
		d_nutrients_to_add.Trans_Fat,
		d_nutrients_to_add.Saturated_Fat,
		d_nutrients_to_add.Sugars,
		d_nutrients_to_add.Fiber,
		d_nutrients_to_add.Sodium,
		d_nutrients_to_add.Iron,
		d_nutrients_to_add.Calcium,
		d_nutrients_to_add.ID,
	)
	if err != nil {
		log.Println("save_daily_nutrients | Error: ", err.Error())
		return err
	}
	return nil
}
func save_intake_ingredient(txn *sql.Tx, intake *models.Intake) error {
	// _, err = txn.Exec(
	// 	`UPDATE account_game_stat SET coins = coins + $1, xp = xp + $2 WHERE account_id = $3`,
	// 	coins, xp, d_nutrients_to_add.Account_Id,
	// )
	// if err != nil {
	// 	log.Println("save_intake_d_nutrients_and_gamestat (update account_game_stat)| Error: ", err.Error())
	// 	return err
	// }
	row := txn.QueryRow(
		`INSERT INTO intake (account_id, date_created, ingredient_mapping_id, amount,	amount_unit, amount_unit_desc, serving_size)
		VALUES ($1, $2, $3, $4, $5, $6, $7)  RETURNING id`,
		intake.Account_Id,
		intake.Date_Created,
		intake.Ingredient_Mapping_Id,
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
	return nil
}

// func save_intake_d_nutrients_and_gamestat(db *sql.DB, d_nutrients_to_add *models.Daily_Nutrients, coins int, xp int, intake *models.Intake) error {
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
// 		log.Println("save_intake_d_nutrients_and_gamestat (update d_nutrients) | Error: ", err.Error())
// 		return err
// 	}
// 	_, err = txn.Exec(
// 		`UPDATE account_game_stat SET coins = coins + $1, xp = xp + $2 WHERE account_id = $3`,
// 		coins, xp, d_nutrients_to_add.Account_Id,
// 	)
// 	if err != nil {
// 		log.Println("save_intake_d_nutrients_and_gamestat (update account_game_stat)| Error: ", err.Error())
// 		return err
// 	}
// 	if intake.Food_Id != 0 && intake.Ingredient_Mapping_Id == 0 {
// 		row := txn.QueryRow(
// 			`INSERT INTO intake (account_id, date_created, food_id, amount,	amount_unit, amount_unit_desc, serving_size)
// 			VALUES ($1, $2, $3, $4, $5, $6, $7)  RETURNING id`,
// 			intake.Account_Id,
// 			intake.Date_Created,
// 			intake.Food_Id,
// 			intake.Amount,
// 			intake.Amount_Unit,
// 			intake.Amount_Unit_Desc,
// 			intake.Serving_Size,
// 		)
// 		err := row.Scan(&intake.ID)
// 		if err != nil {
// 			log.Println("save_intake_d_nutrients_and_gamestat (insert intake)| Error: ", err.Error())
// 			return err
// 		}
// 	}
// 	if intake.Food_Id == 0 && intake.Ingredient_Mapping_Id != 0 {
// 		row := txn.QueryRow(
// 			`INSERT INTO intake (account_id, date_created, recipe_id, amount,	amount_unit, amount_unit_desc, serving_size)
// 			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
// 			intake.Account_Id,
// 			intake.Date_Created,
// 			intake.Amount,
// 			intake.Amount_Unit,
// 			intake.Amount_Unit_Desc,
// 			intake.Serving_Size,
// 		)
// 		err := row.Scan(&intake.ID)
// 		if err != nil {
// 			log.Println("save_intake_d_nutrients_and_gamestat (insert intake)| Error: ", err.Error())
// 			return err
// 		}
// 	}
// 	err = txn.Commit()
// 	if err != nil {
// 		txn.Rollback()
// 		log.Println("save_intake_d_nutrients_and_gamestat (commit) | Error: ", err.Error())
// 		return err
// 	}
// 	return nil
// }
