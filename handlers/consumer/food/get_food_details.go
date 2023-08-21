package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/food"
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
	food := models.Food{}
	nutrient := models.Nutrient{}
	// querying food
	row := query_food_and_nutrient(reqData.Food_ID, reqData.Barcode, db)
	// scanning food
	if row == nil {
		log.Println("Get_Food_Details | Error on query validation: Food_ID and Barcode empty ")
		return utilities.Send_Error(c, "Invalid data sent", fiber.StatusBadRequest)
	}
	err = scan_food_and_nutrient(row, &food, &nutrient)
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
	response.Food = food
	response.Nutrient = nutrient
	return c.Status(fiber.StatusOK).JSON(response)
}

func query_food_and_nutrient(food_id uint, barcode string, db *sql.DB) *sql.Row {
	if food_id != 0 {
		row := db.QueryRow(`SELECT
				food.id,
				food.name,
				food.name_ph,
				food.name_owner,
				food.barcode,
				coalesce(food.thumbnail_image_link, ''),
				coalesce(food.food_desc, ''),
				coalesce(food.category_id, 0),
				food.food_type_id,
				food.owner_id,
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
	if barcode != "" {
		row := db.QueryRow(`SELECT
				food.id,
				food.name,
				food.name_ph,
				food.name_owner,
				food.barcode,
				coalesce(food.thumbnail_image_link, ''),
				coalesce(food.food_desc, ''),
				coalesce(food.category_id, 0),
				food.food_type_id,
				food.owner_id,
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
			WHERE food.barcode = $1`,
			barcode,
		)
		return row
	}
	return nil
}
func scan_food_and_nutrient(row *sql.Row, food *models.Food, nutrient *models.Nutrient) error {
	if err := row.
		Scan(
			&food.ID,
			&food.Name,
			&food.Name_Ph,
			&food.Name_Owner,
			&food.Barcode,
			&food.Thumbnail_Image_Link,
			&food.Food_Desc,
			&food.Category_Id,
			&food.Food_Type_Id,
			&food.Owner_Id,
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
