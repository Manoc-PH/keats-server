package handlers

import (
	"database/sql"
	"log"
	"net/url"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/setup"
	"server/utilities"
	"strconv"

	cld "github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Post_Recipe(c *fiber.Ctx, db *sql.DB) error {
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

	// Validating ingredients
	for _, item := range reqData.Recipe_Ingredients {
		if item.Ingredient_Mapping_Id != 0 {
			exists := ingredient_exists(item.Ingredient_Mapping_Id, db)
			if exists == false {
				log.Println("Post_Recipe | Error: user sent an ingredient that doesn't exist. Ingredient_Mapping_Id: ", item.Ingredient_Mapping_Id)
				return utilities.Send_Error(c, "An ingredient does not exist.", fiber.StatusBadRequest)
			}
		}
		if item.Food_Id != 0 {
			exists := food_exists(item.Food_Id, db)
			if exists == false {
				log.Println("Post_Recipe | Error: user sent an ingredient that doesn't exist. Food_Id: ", item.Food_Id)
				return utilities.Send_Error(c, "An ingredient does not exist.", fiber.StatusBadRequest)
			}
		}
	}

	// Saving Recipe
	err = save_recipe_txn(reqData, db, owner_id)
	if err != nil {
		return utilities.Send_Error(c, "An error occured in saving recipe", fiber.StatusInternalServerError)
	}

	// Generating signature
	strTimestamp := strconv.FormatInt(reqData.Timestamp.Unix(), 10)
	signature, err := cld.SignParameters(url.Values{"timestamp": []string{strTimestamp}}, setup.CloudinaryConfig.APISecret)
	response := schemas.Res_Post_Recipe{
		Recipe:              reqData.Recipe,
		Recipe_Ingredients:  reqData.Recipe_Ingredients,
		Recipe_Instructions: reqData.Recipe_Instructions,
		Signature:           signature,
		Timestamp:           strTimestamp,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
func ingredient_exists(ingredient_mapping_id uint, db *sql.DB) bool {
	row := db.QueryRow(`SELECT id FROM ingredient_mapping WHERE id = $1`, ingredient_mapping_id)
	found_id := 0
	err := row.Scan(&found_id)
	if err != nil {
		return false
	}
	if found_id != 0 {
		return true
	}
	return false
}
func food_exists(food_id uint, db *sql.DB) bool {
	row := db.QueryRow(`SELECT id FROM food WHERE id = $1`, food_id)
	found_id := 0
	err := row.Scan(&found_id)
	if err != nil {
		return false
	}
	if found_id != 0 {
		return true
	}
	return false
}
func save_recipe_txn(recipe *schemas.Req_Post_Recipe, db *sql.DB, owner_id uuid.UUID) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("Post_Recipe | Error on save_recipe_txn: ", err.Error())
		return err
	}
	// Saving Recipe
	err = save_recipe_details(txn, recipe, owner_id)
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

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println(" (commit) | Error: ", err.Error())
		return err
	}
	return nil
}
func save_recipe_details(txn *sql.Tx, recipe *schemas.Req_Post_Recipe, owner_id uuid.UUID) error {
	// Saving Recipe Details
	// TODO SAVE CATEGORY ID
	row := txn.QueryRow(`
	INSERT INTO recipe(
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
		description)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING id
`,
		recipe.Recipe.Name,
		recipe.Recipe.Name_Ph,
		recipe.Recipe.Name_Owner,
		owner_id,
		recipe.Recipe.Date_Created,
		// recipe.Recipe.Category_Id,
		recipe.Recipe.Thumbnail_Image_Link,
		recipe.Recipe.Main_Image_Link,
		recipe.Recipe.Likes,
		recipe.Recipe.Rating,
		recipe.Recipe.Servings,
		recipe.Recipe.Servings_Size,
		recipe.Recipe.Prep_Time,
		recipe.Recipe.Description)
	err := row.Scan(&recipe.Recipe.ID)
	if err != nil {
		log.Println("Post_Recipe | Error on save_recipe_details: ", err.Error())
		return err
	}
	return nil
}
func save_recipe_ingredients(txn *sql.Tx, recipe *schemas.Req_Post_Recipe) error {
	// Saving Recipe Ingredients
	stmt, err := txn.Prepare(
		`INSERT INTO recipe_ingredient (
				food_id,
				ingredient_mapping_id,
				amount,
				amount_unit,
				amount_unit_desc,
				serving_size)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
	)
	if err != nil {
		log.Println("save_recipe_ingredient (Prepare) | Error: ", err.Error())
		return err
	}
	defer stmt.Close()
	// Insert each row
	for i, item := range recipe.Recipe_Ingredients {
		row := stmt.QueryRow(
			item.Food_Id,
			item.Ingredient_Mapping_Id,
			item.Amount,
			item.Amount_Unit,
			item.Amount_Unit_Desc,
			item.Serving_Size,
		)
		id := 0
		err = row.Scan(&id)
		recipe.Recipe_Ingredients[i].ID = uint(id)
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
				recipe_id,
				instruction_description,
				step_num)
			VALUES ($1, $2, $3) RETURNING id`,
	)
	if err != nil {
		log.Println("save_recipe_ingredient (Prepare) | Error: ", err.Error())
		return err
	}
	defer stmt.Close()
	// Insert each row
	for i, item := range recipe.Recipe_Instructions {
		row := stmt.QueryRow(
			recipe.Recipe.ID,
			item.Instruction_Description,
			item.Step_Num,
		)
		id := 0
		err = row.Scan(&id)
		recipe.Recipe_Instructions[i].ID = uint(id)
		if err != nil {
			log.Println("save_recipe_ingredient (Exec) | Error: ", err.Error())
			return err
		}
	}
	return nil
}
