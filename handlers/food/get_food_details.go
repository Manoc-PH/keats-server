package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/food"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Get_Food_Details(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Food_Details | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Food_Details)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Food_Details | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	response := schemas.Res_Get_Food_Details{}
	// querying food
	row := query_food(db, reqData.Food_ID)
	// scanning food
	err = scan_food(row, &response)
	if err != nil && err == sql.ErrNoRows {
		log.Println("Get_Food_Details | error in scanning food: ", err.Error())
		return utilities.Send_Error(c, "Food does not exist", fiber.StatusBadRequest)
	}
	// Server Error
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get_Food_Details | error in scanning food: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	// querying food
	images, err := query_and_scan_food_images(db, reqData.Food_ID)
	if err != nil {
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	response.Food_Images = images
	return c.Status(fiber.StatusOK).JSON(response)
}

func query_food(db *sql.DB, food_id uint) *sql.Row {
	row := db.QueryRow(`SELECT
			food.id, 
			food.name,
			food.name_ph,
			food.name_brand,
			food.thumbnail_image_link,
			food.food_desc,
			food.food_nutrient_id,
			food.food_brand_type_id,
			food.food_category_id,
			food.food_brand_id,
			--	FOOD NUTRIENT
			food_nutrient.amount,
			food_nutrient.amount_unit,
			food_nutrient.amount_unit_desc,
			food_nutrient.serving_size,
			food_nutrient.calories,
			food_nutrient.protein,
			food_nutrient.carbs,
			food_nutrient.fats,
			food_nutrient.trans_fat,
			food_nutrient.saturated_fat,
			food_nutrient.sugars,
			food_nutrient.sodium,
			-- FOOD BRAND TYPE
			food_brand_type.name,
			food_brand_type.brand_type_desc,
			-- FOOD BRAND
			food_brand.id
			-- 
		FROM food
		LEFT JOIN food_nutrient ON food.food_nutrient_id = food_nutrient.id
		LEFT JOIN food_brand_type ON food.food_brand_type_id = food_brand_type.id
		LEFT JOIN food_brand ON food.food_brand_id = food_brand.id
		WHERE food.id = $1`, food_id,
	)
	return row
}
func scan_food(row *sql.Row, food *schemas.Res_Get_Food_Details) error {
	err := row.Scan(
		&food.ID,
		&food.Name,
		&food.Name_Ph,
		&food.Name_Brand,
		&food.Thumbnail_Image_Link,
		&food.Food_Desc,
		&food.Food_Nutrient_Id,
		&food.Food_Brand_Type_Id,
		&food.Food_Category_Id,
		&food.Food_Brand_Id,
		// FOOD NUTRIENT
		&food.Food_Nutrients.Amount,
		&food.Food_Nutrients.Amount_Unit,
		&food.Food_Nutrients.Amount_Unit_Desc,
		&food.Food_Nutrients.Serving_Size,
		&food.Food_Nutrients.Calories,
		&food.Food_Nutrients.Protein,
		&food.Food_Nutrients.Carbs,
		&food.Food_Nutrients.Fats,
		&food.Food_Nutrients.Trans_Fat,
		&food.Food_Nutrients.Saturated_Fat,
		&food.Food_Nutrients.Sugars,
		&food.Food_Nutrients.Sodium,
		// FOOD BRAND TYPE
		&food.Food_Brand_Type.Name,
		&food.Food_Brand_Type.Brand_Type_Desc,
		// FOOD BRAND
		&food.Food_Brand.ID,
	)
	return err
}

func query_and_scan_food_images(db *sql.DB, food_id uint) ([]models.Food_Image, error) {
	rows, err := db.Query(`SELECT
			id,
			food_id,
			name_file,
			amount,
			amount_unit,
			amount_unit_desc
		FROM food_image
		WHERE food_id = $1`, food_id,
	)
	if err != nil {
		log.Println("Get_Food_Details | error in querying food: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	images := make([]models.Food_Image, 0, 10)
	for rows.Next() {
		var new_image = models.Food_Image{}
		if err := rows.
			Scan(
				&new_image.ID,
				&new_image.Food_Id,
				&new_image.Name_File,
				&new_image.Amount,
				&new_image.Amount_Unit,
				&new_image.Amount_Unit_Desc,
			); err != nil {
			log.Println("Get_Food_Details | error in scanning image: ", err.Error())
			return nil, err
		}
		images = append(images, new_image)
	}
	return images, err
}
