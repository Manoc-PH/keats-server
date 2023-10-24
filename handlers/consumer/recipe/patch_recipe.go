package handlers

import (
	"database/sql"
	"log"
	"server/constants"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Patch_Recipe(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Patch_Recipe | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Patch_Recipe)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Patch_Recipe | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// DB transaction
	txn, err := db.Begin()
	if err != nil {
		log.Println("Patch_Recipe | Error on db.Begin(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// Updating Recipe details
	err = update_recipe_details(txn, &reqData.Recipe)
	if err != nil {
		log.Println("Patch_Recipe | Error on update_recipe_details: ", err.Error())
		return utilities.Send_Error(c, "An error occured in updating recipe details", fiber.StatusInternalServerError)
	}

	// Updating Recipe Ingredients
	if len(reqData.Recipe_Ingredients) > 0 {
		err = update_recipe_ingredients(txn, &reqData.Recipe_Ingredients, reqData.Recipe)
		if err != nil {
			log.Println("Patch_Recipe | Error on update_recipe_ingredients: ", err.Error())
			return utilities.Send_Error(c, "An error occured in updating recipe details", fiber.StatusInternalServerError)
		}
	}

	// Updating Recipe Instructions
	if len(reqData.Recipe_Instructions) > 0 {
		err = update_instructions(txn, &reqData.Recipe_Instructions, reqData.Recipe)
		if err != nil {
			log.Println("Patch_Recipe | Error on update_instructions: ", err.Error())
			return utilities.Send_Error(c, "An error occured in updating recipe details", fiber.StatusInternalServerError)
		}
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("Patch_Recipe | Error on txn.Commit(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}

func update_recipe_details(tx *sql.Tx, data *schemas.Recipe_Patch) error {
	// TODO ADD category_id, thumbnail_image_link, main_image_link
	_, err := tx.Exec(`UPDATE recipe SET 
			name = $1,
			name_ph = $2,
			servings = $3,
			servings_size = $4,
			prep_time = $5,
			description = $6
		WHERE id = $7`,
		data.Name,
		data.Name_Ph,
		data.Servings,
		data.Servings_Size,
		data.Prep_Time,
		data.Description,
		data.ID,
	)

	return err
}
func update_recipe_ingredients(tx *sql.Tx, data *[]schemas.Recipe_Patch_Ingredient, recipe schemas.Recipe_Patch) error {
	stmtInsert, err := tx.Prepare(`INSERT INTO recipe_ingredient (
			food_id,
			ingredient_mapping_id,
			amount,
			amount_unit,
			amount_unit_desc,
			serving_size,
			recipe_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
	)
	if err != nil {
		log.Println(" Error on update_recipe_ingredients")
		return nil
	}
	stmtUpdate, err := tx.Prepare(`UPDATE recipe_ingredient SET
			food_id = $1,
			ingredient_mapping_id = $2,
			amount = $3,
			amount_unit = $4,
			amount_unit_desc = $5,
			serving_size = $6
		WHERE id = $7`)
	if err != nil {
		log.Println(" Error on update_recipe_ingredients")
		return nil
	}
	stmtDelete, err := tx.Prepare(`DELETE FROM recipe_ingredient WHERE id = $1`)
	if err != nil {
		log.Println(" Error on update_recipe_ingredients")
		return nil
	}
	for _, item := range *data {
		if item.Action_Type == constants.Action_Types.Insert {
			_, err = stmtInsert.Exec(
				item.Food_Id,
				item.Ingredient_Mapping_Id,
				item.Amount,
				item.Amount_Unit,
				item.Amount_Unit_Desc,
				item.Serving_Size,
				recipe.ID)
			if err != nil {
				return err
			}
		}
		if item.Action_Type == constants.Action_Types.Update {
			_, err = stmtUpdate.Exec(
				item.Food_Id,
				item.Ingredient_Mapping_Id,
				item.Amount,
				item.Amount_Unit,
				item.Amount_Unit_Desc,
				item.Serving_Size,
				item.ID)
			if err != nil {
				return err
			}
		}
		if item.Action_Type == constants.Action_Types.Delete {
			_, err = stmtDelete.Exec(item.ID)
			if err != nil {
				return err
			}
		}
	}

	// getting the new nurtients
	rows, err := tx.Query(`SELECT 
			id,
			food_id,
			ingredient_mapping_id,
			amount,
			amount_unit,
			amount_unit_desc,
			serving_size,
			recipe_id
		FROM recipe_ingredient 
		WHERE recipe_id = $1`, recipe.ID)
	defer rows.Close()

	if err != nil {
		log.Println(" Error on update_recipe_ingredients")
		return err
	}

	newIngredients := []schemas.Recipe_Ingredient_Schema{}
	for rows.Next() {
		newIngredient := schemas.Recipe_Ingredient_Schema{}
		err := rows.Scan(
			&newIngredient.ID,
			&newIngredient.Food_Id,
			&newIngredient.Ingredient_Mapping_Id,
			&newIngredient.Amount,
			&newIngredient.Amount_Unit,
			&newIngredient.Amount_Unit_Desc,
			&newIngredient.Serving_Size,
			&newIngredient.Recipe_Id,
		)
		if err != nil {
			log.Println(" Error on update_recipe_ingredients")
			return err
		}
		newIngredients = append(newIngredients, newIngredient)
	}

	nutrients, err := generate_nutrients_tx(tx, &newIngredients, recipe.Servings)
	nutrients.ID = recipe.Nutrient_Id
	if err != nil {
		log.Println(" Error on update_recipe_ingredients")
		return err
	}

	// Updating Nutrient
	err = update_nutrient(tx, nutrients)
	return err
}
func update_nutrient(tx *sql.Tx, nutrient *models.Nutrient) error {
	_, err := tx.Exec(`UPDATE nutrient SET
			amount = $1,
			amount_unit = $2,
			amount_unit_desc = $3,
			serving_size = $4,
			serving_total = $5,
			calories = $6,
			protein = $7,
			carbs = $8,
			fats = $9,
			trans_fat = $10,
			saturated_fat = $11,
			sugars = $12,
			fiber = $13,
			sodium = $14,
			iron = $15,
			calcium = $16
		WHERE id = $17`,
		nutrient.Amount,
		nutrient.Amount_Unit,
		nutrient.Amount_Unit_Desc,
		nutrient.Serving_Size,
		nutrient.Serving_Total,
		nutrient.Calories,
		nutrient.Protein,
		nutrient.Carbs,
		nutrient.Fats,
		nutrient.Trans_Fat,
		nutrient.Saturated_Fat,
		nutrient.Sugars,
		nutrient.Fiber,
		nutrient.Sodium,
		nutrient.Iron,
		nutrient.Calcium,
		nutrient.ID,
	)
	if err != nil {
		log.Println(" Error on update_nutrient")
		return err
	}
	return nil
}
func update_instructions(tx *sql.Tx, data *[]schemas.Recipe_Patch_Instruction, recipe schemas.Recipe_Patch) error {
	stmtInsert, err := tx.Prepare(`INSERT INTO recipe_instruction (
			recipe_id,
			instruction_description,
			step_num)
		VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		log.Println("Error on update_instructions(stmtInsert)")
		return nil
	}
	stmtUpdate, err := tx.Prepare(`UPDATE recipe_instruction SET
		instruction_description = $1,
		step_num = $2
	WHERE id = $3`)
	if err != nil {
		log.Println("Error on update_instructions(stmtUpdate)")
		return nil
	}
	stmtDelete, err := tx.Prepare(`DELETE FROM recipe_instruction WHERE id = $1`)
	if err != nil {
		log.Println("Error on update_instructions(stmtDelete)")
		return nil
	}
	for _, item := range *data {
		if item.Action_Type == constants.Action_Types.Insert {
			_, err = stmtInsert.Exec(
				recipe.ID,
				item.Instruction_Description,
				item.Step_Num)
			if err != nil {
				log.Println(" Error on update_instructions(stmtInsert.Exec)")
				return err
			}
		}
		if item.Action_Type == constants.Action_Types.Update {
			_, err = stmtUpdate.Exec(
				item.Instruction_Description,
				item.Step_Num,
				item.ID)
			if err != nil {
				log.Println(" Error on update_instructions(stmtUpdate.Exec)")
				return err
			}
		}
		if item.Action_Type == constants.Action_Types.Delete {
			_, err = stmtDelete.Exec(item.ID)
			if err != nil {
				log.Println(" Error on update_instructions(stmtDelete.Exec)")
				return err
			}
		}
	}
	return nil
}

func get_ingredient_nutrient_tx(ingredient_mapping_id uint, tx *sql.Tx, nutrient *models.Nutrient) error {
	row := tx.QueryRow(`SELECT
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
		JOIN nutrient ON ingredient_mapping.nutrient_id = nutrient.id
		WHERE ingredient_mapping.id = $1`,
		ingredient_mapping_id,
	)
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
		log.Println("Error on get_ingredient_nutrient_tx(Scan)")
		return err
	}
	return nil
}
func get_food_nutrient_tx(food_id uint, tx *sql.Tx, nutrient *models.Nutrient) error {
	row := tx.QueryRow(`SELECT 
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
		log.Println("Error on get_food_nutrient_tx(Scan)")
		return err
	}
	return nil
}
func generate_nutrients_tx(tx *sql.Tx, reqData *[]schemas.Recipe_Ingredient_Schema, servings uint) (*models.Nutrient, error) {
	nutrient := new(models.Nutrient)
	for _, item := range *reqData {
		item_nutrient := new(models.Nutrient)
		if item.Ingredient_Mapping_Id != 0 {
			err := get_ingredient_nutrient_tx(item.Ingredient_Mapping_Id, tx, item_nutrient)
			if err != nil {
				log.Println("Error on get_ingredient_nutrient_tx")
				return nil, err
			}
		}
		if item.Food_Id != 0 {
			err := get_food_nutrient_tx(item.Food_Id, tx, item_nutrient)
			if err != nil {
				log.Println("Error on get_food_nutrient_tx")
				return nil, err
			}
		}
		multiplier := item.Amount * 0.01
		nutrient.Serving_Total += item.Amount
		nutrient.Calories += ((item_nutrient.Calories) * multiplier)
		nutrient.Protein += ((item_nutrient.Protein) * multiplier)
		nutrient.Carbs += ((item_nutrient.Carbs) * multiplier)
		nutrient.Fats += ((item_nutrient.Fats) * multiplier)
		nutrient.Trans_Fat += ((item_nutrient.Trans_Fat) * multiplier)
		nutrient.Saturated_Fat += ((item_nutrient.Saturated_Fat) * multiplier)
		nutrient.Sugars += ((item_nutrient.Sugars) * multiplier)
		nutrient.Fiber += ((item_nutrient.Fiber) * multiplier)
		nutrient.Sodium += ((item_nutrient.Sodium) * multiplier)
		nutrient.Iron += ((item_nutrient.Iron) * multiplier)
		nutrient.Calcium += ((item_nutrient.Calcium) * multiplier)
	}
	// calculating nutrients
	divider := nutrient.Serving_Total * 0.01
	nutrient.Amount = 100
	nutrient.Amount_Unit = "g"
	nutrient.Amount_Unit_Desc = "grams"
	nutrient.Serving_Size = nutrient.Serving_Total / float32(servings)
	nutrient.Calories = nutrient.Calories / divider
	nutrient.Protein = nutrient.Protein / divider
	nutrient.Carbs = nutrient.Carbs / divider
	nutrient.Fats = nutrient.Fats / divider
	nutrient.Trans_Fat = nutrient.Trans_Fat / divider
	nutrient.Saturated_Fat = nutrient.Saturated_Fat / divider
	nutrient.Sugars = nutrient.Sugars / divider
	nutrient.Fiber = nutrient.Fiber / divider
	nutrient.Sodium = nutrient.Sodium / divider
	nutrient.Iron = nutrient.Iron / divider
	nutrient.Calcium = nutrient.Calcium / divider

	return nutrient, nil
}
