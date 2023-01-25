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
			food.food_nutrient_id,
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
			food_nutrient.sodium
		FROM food
		LEFT JOIN food_nutrient ON food.food_nutrient_id = food_nutrient.id
		WHERE food.id = $1`, food_id,
	)
	return row
}
func scan_food(row *sql.Row, food *schemas.Res_Get_Food_Details) error {
	err := row.Scan(
		&food.Food_Details.ID,
		&food.Food_Details.Name,
		&food.Food_Details.Name_Ph,
		&food.Food_Details.Name_Brand,
		&food.Food_Details.Food_Nutrient_Id,
		&food.Food_Details.Amount,
		&food.Food_Details.Amount_Unit,
		&food.Food_Details.Amount_Unit_Desc,
		&food.Food_Details.Serving_Size,
		&food.Food_Details.Calories,
		&food.Food_Details.Protein,
		&food.Food_Details.Carbs,
		&food.Food_Details.Fats,
		&food.Food_Details.Trans_Fat,
		&food.Food_Details.Saturated_Fat,
		&food.Food_Details.Sugars,
		&food.Food_Details.Sodium,
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
