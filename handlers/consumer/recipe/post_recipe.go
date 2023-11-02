package handlers

import (
	"database/sql"
	"log"
	"net/url"
	"server/constants"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/recipe"
	"server/setup"
	"server/utilities"
	"strconv"
	"time"

	cld "github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
)

// TODO UPDATE THE HANDLER FOR IMAGE UPLOAD. THE URL SHOULD BE GENERATED HERE
func Post_Recipe(c *fiber.Ctx, db *sql.DB, db_search *meilisearch.Client) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Recipe | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Post_Recipe)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Recipe | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// Assigning id to recipe
	reqData.Recipe.ID = uuid.New()

	// Generating the nutrients
	nutrient, err := generate_nutrients(db, &reqData.Recipe_Ingredients, reqData.Recipe.Servings)
	if err != nil {
		log.Println("Post_Recipe | Error on generate_nutrients: ", err.Error())
		return utilities.Send_Error(c, "An ingredient does not exist.", fiber.StatusBadRequest)
	}
	nutrient.Parent_ID = reqData.Recipe.ID

	// Saving Recipe
	err = save_recipe_txn(reqData, db, db_search, owner_id, nutrient)
	if err != nil {
		return utilities.Send_Error(c, "An error occured in saving recipe", fiber.StatusInternalServerError)
	}

	// Generating signature
	strTimestamp := strconv.FormatInt(reqData.Timestamp.Unix(), 10)
	signature, err := cld.SignParameters(url.Values{"timestamp": []string{strTimestamp}}, setup.Cloudinary_Config.APISecret)
	response := schemas.Res_Post_Recipe{
		Recipe:              reqData.Recipe,
		Recipe_Ingredients:  reqData.Recipe_Ingredients,
		Recipe_Instructions: reqData.Recipe_Instructions,
		Nutrient:            *nutrient,
		Signature:           signature,
		Timestamp:           strTimestamp,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
func get_ingredient_nutrient(ingredient_mapping_id uuid.UUID, db *sql.DB, nutrient *models.Nutrient) error {
	row := db.QueryRow(`SELECT
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
		return err
	}
	return nil
}
func get_food_nutrient(food_id uuid.UUID, db *sql.DB, nutrient *models.Nutrient) error {
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
		return err
	}
	return nil
}
func generate_nutrients(db *sql.DB, reqData *[]schemas.Recipe_Ingredient_Schema, servings uint) (*models.Nutrient, error) {
	nutrient := new(models.Nutrient)
	for _, item := range *reqData {
		item_nutrient := new(models.Nutrient)
		if item.Ingredient_Mapping_Id != constants.Empty_UUID {
			err := get_ingredient_nutrient(item.Ingredient_Mapping_Id, db, item_nutrient)
			if err != nil {
				return nil, err
			}
		}
		if item.Food_Id != constants.Empty_UUID {
			err := get_food_nutrient(item.Food_Id, db, item_nutrient)
			if err != nil {
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
	nutrient.ID = uuid.New()
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
func save_recipe_txn(
	recipe *schemas.Req_Post_Recipe,
	db *sql.DB,
	db_search *meilisearch.Client,
	owner_id uuid.UUID,
	nutrient *models.Nutrient,
) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("Post_Recipe | Error on save_recipe_txn: ", err.Error())
		return err
	}
	// Saving Nutrient
	err = save_nutrient(txn, nutrient)
	if err != nil {
		log.Println("Post_Recipe | Error on save_nutrient: ", err.Error())
		return err
	}
	// Saving Recipe
	err = save_recipe_details(txn, recipe, nutrient, owner_id)
	if err != nil {
		log.Println("Post_Recipe | Error on save_recipe_details: ", err.Error())
		return err
	}
	// Saving Recipe Ingredients
	err = save_recipe_ingredients(txn, recipe)
	if err != nil {
		log.Println("Post_Recipe | Error on save_recipe_ingredients: ", err.Error())
		return err
	}
	// Saving Recipe Instructions
	err = save_recipe_instructions(txn, recipe)
	if err != nil {
		log.Println("Post_Recipe | Error on save_recipe_instructions: ", err.Error())
		return err
	}
	// Saving Recipe to Meilisearch
	err = save_recipe_to_meili(db_search, recipe)
	if err != nil {
		log.Println("Post_Recipe | Error on save_recipe_to_meili: ", err.Error())
		txn.Rollback()
		return err
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println(" (commit) | Error: ", err.Error())
		return err
	}
	return nil
}
func save_recipe_details(txn *sql.Tx, recipe *schemas.Req_Post_Recipe, nutrient *models.Nutrient, owner_id uuid.UUID) error {
	recipe.Recipe.Servings_Size = nutrient.Serving_Size
	// Saving Recipe Details
	// TODO SAVE CATEGORY ID
	_, err := txn.Exec(`
		INSERT INTO recipe(
			id,
			name,
			name_ph,
			name_owner,
			owner_id,
			date_created,
			-- category_id,
			thumbnail_image_link,
			main_image_link,
			likes,
			rating,
			servings,
			servings_size,
			prep_time,
			description,
			nutrient_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		recipe.Recipe.ID,
		recipe.Recipe.Name,
		recipe.Recipe.Name_Ph,
		recipe.Recipe.Name_Owner,
		owner_id,
		time.Now(),
		// recipe.Recipe.Category_Id,
		recipe.Recipe.Thumbnail_Image_Link,
		recipe.Recipe.Main_Image_Link,
		recipe.Recipe.Likes,
		recipe.Recipe.Rating,
		recipe.Recipe.Servings,
		nutrient.Serving_Size,
		recipe.Recipe.Prep_Time,
		recipe.Recipe.Description,
		nutrient.ID)
	if err != nil {
		log.Println("Post_Recipe | Error on save_recipe_details: ", err.Error())
		return err
	}
	return nil
}
func save_recipe_ingredients(txn *sql.Tx, recipe *schemas.Req_Post_Recipe) error {
	// Saving Recipe Ingredients
	stmt_ingredient, err := txn.Prepare(
		`INSERT INTO recipe_ingredient (
				id,
				ingredient_mapping_id,
				amount,
				amount_unit,
				amount_unit_desc,
				serving_size,
				recipe_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
	)
	if err != nil {
		log.Println("save_recipe_ingredient (Prepare 1) | Error: ", err.Error())
		return err
	}
	defer stmt_ingredient.Close()
	stmt_food, err := txn.Prepare(
		`INSERT INTO recipe_ingredient (
				id,
				food_id,
				amount,
				amount_unit,
				amount_unit_desc,
				serving_size,
				recipe_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
	)
	if err != nil {
		log.Println("save_recipe_ingredient (Prepare 2) | Error: ", err.Error())
		return err
	}
	defer stmt_food.Close()
	// Insert each row
	for i, item := range recipe.Recipe_Ingredients {
		id := uuid.New()
		recipe.Recipe_Ingredients[i].ID = id
		var err error
		if item.Food_Id != constants.Empty_UUID {
			_, err = stmt_food.Exec(
				id,
				item.Food_Id,
				item.Amount,
				item.Amount_Unit,
				item.Amount_Unit_Desc,
				item.Serving_Size,
				recipe.Recipe.ID,
			)
		}
		if item.Ingredient_Mapping_Id != constants.Empty_UUID {
			_, err = stmt_ingredient.Exec(
				id,
				item.Ingredient_Mapping_Id,
				item.Amount,
				item.Amount_Unit,
				item.Amount_Unit_Desc,
				item.Serving_Size,
				recipe.Recipe.ID,
			)
		}
		if err != nil {
			log.Println("save_recipe_ingredient (Exec) | Error: ", err.Error())
			return err
		}
	}
	return nil
}
func save_recipe_instructions(txn *sql.Tx, recipe *schemas.Req_Post_Recipe) error {
	// Saving Recipe Instructions
	stmt, err := txn.Prepare(
		`INSERT INTO recipe_instruction (
				id,
				recipe_id,
				instruction_description,
				step_num)
			VALUES ($1, $2, $3, $4)`,
	)
	if err != nil {
		log.Println("save_recipe_ingredient (Prepare) | Error: ", err.Error())
		return err
	}
	defer stmt.Close()
	// Insert each row
	for i, item := range recipe.Recipe_Instructions {
		id := uuid.New()
		recipe.Recipe_Instructions[i].ID = id
		_, err := stmt.Exec(
			id,
			recipe.Recipe.ID,
			item.Instruction_Description,
			item.Step_Num,
		)
		if err != nil {
			log.Println("save_recipe_ingredient (Exec) | Error: ", err.Error())
			return err
		}
	}
	return nil
}
func save_nutrient(txn *sql.Tx, nutrient *models.Nutrient) error {
	_, err := txn.Exec(`INSERT INTO nutrient
			(id,
			parent_id,
			amount,
			amount_unit,
			amount_unit_desc,
			serving_size,
			serving_total,
			calories,
			protein,
			carbs,
			fats,
			trans_fat,
			saturated_fat,
			sugars,
			fiber,
			sodium,
			iron,
			calcium)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		nutrient.ID,
		nutrient.Parent_ID,
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
	)
	if err != nil {
		return err
	}
	return nil
}
func save_recipe_to_meili(db_search *meilisearch.Client, recipe *schemas.Req_Post_Recipe) error {
	new_item := map[string]interface{}{
		"id":                   recipe.Recipe.ID,
		"name":                 recipe.Recipe.Name,
		"name_ph":              recipe.Recipe.Name_Ph,
		"name_owner":           recipe.Recipe.Name_Owner,
		"thumbnail_image_link": recipe.Recipe.Thumbnail_Image_Link,
		"main_image_link":      recipe.Recipe.Main_Image_Link,
		"rating":               recipe.Recipe.Rating,
		"rating_count":         recipe.Recipe.Rating_Count,
	}
	_, err := db_search.Index("recipes").AddDocuments(new_item, "id")
	if err != nil {
		return err
	}
	return nil
}
