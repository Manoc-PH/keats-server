package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Get_Recipe_Details(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("get_recipe_details | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Recipe_Details)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("get_recipe_details | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	response := new(schemas.Res_Get_Recipe_Details)
	// * getting recipe details
	err = get_recipe_details_and_nutrients(db, reqData.Recipe_Id, &response.Recipe, &response.Nutrients)
	if err != nil {
		log.Println("get_recipe_details | Error on get_recipe_details_and_nutrients: ", err.Error())
		return utilities.Send_Error(c, "An error occured in fetching recipe", fiber.StatusInternalServerError)
	}
	// * getting recipe images
	err = get_recipe_images(db, reqData.Recipe_Id, &response.Recipe_Images)
	if err != nil {
		log.Println("get_recipe_details | Error on get_recipe_images: ", err.Error())
		return utilities.Send_Error(c, "An error occured in fetching recipe", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func get_recipe_details_and_nutrients(db *sql.DB, recipe_id uuid.UUID, recipe *models.Recipe, nutrient *models.Nutrient) error {
	// TODO MAKE SURE CATEGORY ID IS NOT 0
	row := db.QueryRow(`SELECT 
			recipe.id,
			recipe.name,
			recipe.name_ph,
			recipe.name_owner,
			recipe.owner_id,
			recipe.date_created,
			COALESCE(recipe.category_id, 0),
			recipe.thumbnail_image_link,
			recipe.main_image_link,
			recipe.likes,
			recipe.rating,
			recipe.rating_count,
			recipe.servings,
			recipe.servings_size,
			recipe.prep_time,
			recipe.description,
			nutrient.id,
			nutrient.amount,
			nutrient.amount_unit,
			nutrient.amount_unit_desc,
			nutrient.serving_size,
			nutrient.serving_total,
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
		FROM recipe
		JOIN nutrient ON recipe.nutrient_id = nutrient.id
		WHERE recipe.id = $1`, recipe_id)
	err := row.Scan(
		&recipe.ID,
		&recipe.Name,
		&recipe.Name_Ph,
		&recipe.Name_Owner,
		&recipe.Owner_Id,
		&recipe.Date_Created,
		&recipe.Category_Id,
		&recipe.Thumbnail_Image_Link,
		&recipe.Main_Image_Link,
		&recipe.Likes,
		&recipe.Rating,
		&recipe.Rating_Count,
		&recipe.Servings,
		&recipe.Servings_Size,
		&recipe.Prep_Time,
		&recipe.Description,
		&nutrient.ID,
		&nutrient.Amount,
		&nutrient.Amount_Unit,
		&nutrient.Amount_Unit_Desc,
		&nutrient.Serving_Size,
		&nutrient.Serving_Total,
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
	)
	if err != nil {
		return err
	}

	return nil
}
func get_recipe_images(db *sql.DB, recipe_id uuid.UUID, recipe_images *[]models.Recipe_Image) error {
	rows, err := db.Query(`SELECT
			id,
			recipe_id,
			name_file,
			name_url,
			amount,
			amount_unit,
			amount_unit_desc
		FROM recipe_image
		WHERE recipe_id = $1`,
		recipe_id,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var recipe_img = models.Recipe_Image{}
		if err := rows.
			Scan(
				&recipe_img.ID,
				&recipe_img.Recipe_Id,
				&recipe_img.Name_File,
				&recipe_img.Name_URL,
				&recipe_img.Amount,
				&recipe_img.Amount_Unit,
				&recipe_img.Amount_Unit_Desc,
			); err != nil {
			return err
		}
		*recipe_images = append(*recipe_images, recipe_img)
	}
	return nil
}
