package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/tracker"
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
		ingredient_mapping := &schemas.Ingredient_Mapping_Schema{}
		daily_nutrients := models.Daily_Nutrients{Account_Id: owner_id}
		nutrients_to_add := models.Nutrient{}
		// Getting ingredient data
		row := query_ingredient(reqData.Ingredient_Mapping_Id, db)
		err = scan_ingredient(row, ingredient_mapping)
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
			log.Println(" (commit) | Error: ", err.Error())
			return err
		}
		response_data.Added_Daily_Nutrients = nutrients_to_add
		response_data.Ingredient = ingredient_mapping
		response_data.Food = nil
		response_data.Intake = new_intake
	}

	// *Food
	if reqData.Food_Id != 0 {
		food := &models.Food{}
		food_nutrients := models.Nutrient{}
		daily_nutrients := models.Daily_Nutrients{Account_Id: owner_id}
		nutrients_to_add := models.Nutrient{}
		// Getting ingredient data
		row := query_food_and_nutrient(reqData.Food_Id, db)
		err = scan_food_and_nutrient(row, food, &food_nutrients)
		if err != nil {
			log.Println("Post_Intake | Error on scanning food: ", err.Error())
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
		calc_nutrients(&nutrients_to_add, &food_nutrients, reqData.Amount)
		nutrients_to_add.ID = daily_nutrients.ID
		new_intake := models.Intake{
			Account_Id:       owner_id,
			Date_Created:     time.Now(),
			Food_Id:          reqData.Food_Id,
			Amount:           reqData.Amount,
			Amount_Unit:      reqData.Amount_Unit,
			Amount_Unit_Desc: reqData.Amount_Unit_Desc,
			Serving_Size:     reqData.Serving_Size,
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
		err = save_intake_food(txn, &new_intake)
		if err != nil {
			return utilities.Send_Error(c, "An error occured in saving intake", fiber.StatusInternalServerError)
		}
		err = txn.Commit()
		if err != nil {
			txn.Rollback()
			log.Println(" (commit) | Error: ", err.Error())
			return err
		}
		response_data.Added_Daily_Nutrients = nutrients_to_add
		response_data.Food = food
		response_data.Ingredient = nil
		response_data.Intake = new_intake
	}

	return c.Status(fiber.StatusOK).JSON(response_data)
}

func query_ingredient(ingredient_mapping_id uint, db *sql.DB) *sql.Row {
	row := db.QueryRow(`SELECT
			ingredient.id, ingredient.name, coalesce(ingredient.name_ph, ''), ingredient.name_owner,
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
			&ingredient_mapping.Ingredient.Name_Owner,

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
func query_food_and_nutrient(food_id uint, db *sql.DB) *sql.Row {
	row := db.QueryRow(`SELECT
			food.id, food.name, food.name_ph, food.name_owner,
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
		food_id,
	)
	return row
}
func scan_food_and_nutrient(row *sql.Row, food *models.Food, nutrient *models.Nutrient) error {
	if err := row.
		Scan(
			&food.ID,
			&food.Name,
			&food.Name_Ph,
			&food.Name_Owner,
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
	// 	log.Println("save_intake_ingredient (update account_game_stat)| Error: ", err.Error())
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
		log.Println("save_intake_ingredient (insert intake)| Error: ", err.Error())
		return err
	}
	return nil
}
func save_intake_food(txn *sql.Tx, intake *models.Intake) error {
	row := txn.QueryRow(
		`INSERT INTO intake (account_id, date_created, food_id, amount,	amount_unit, amount_unit_desc, serving_size)
		VALUES ($1, $2, $3, $4, $5, $6, $7)  RETURNING id`,
		intake.Account_Id,
		intake.Date_Created,
		intake.Food_Id,
		intake.Amount,
		intake.Amount_Unit,
		intake.Amount_Unit_Desc,
		intake.Serving_Size,
	)
	err := row.Scan(&intake.ID)
	if err != nil {
		log.Println("save_intake_food (insert intake)| Error: ", err.Error())
		return err
	}
	return nil
}
