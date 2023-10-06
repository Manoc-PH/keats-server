package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/tracker"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TODO Update this handler, add handler for food
// Gets the details of the intake
func Get_Intake_Details(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, Owner_Id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Intake_Details | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Intake_Details)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Intake_Details | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	intake := models.Intake{}
	// querying intake
	row := query_intake(db, Owner_Id, reqData.Intake_ID)
	// scanning intake
	err = scan_intake(row, &intake)
	if err != nil && err == sql.ErrNoRows {
		log.Println("Get_Intake_Details | error in scanning intake: ", err.Error())
		return utilities.Send_Error(c, "Intake does not exist", fiber.StatusBadRequest)
	}
	// Server Error
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get_Intake_Details | error in scanning intake: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	response := schemas.Res_Get_Intake_Details{
		ID:                    intake.ID,
		Account_Id:            intake.Account_Id,
		Food_Id:               intake.Food_Id,
		Ingredient_Mapping_Id: intake.Ingredient_Mapping_Id,
		Date_Created:          intake.Date_Created,
		Amount:                intake.Amount,
		Amount_Unit:           intake.Amount_Unit,
		Amount_Unit_Desc:      intake.Amount_Unit_Desc,
		Serving_Size:          intake.Serving_Size,
	}
	if intake.Ingredient_Mapping_Id != 0 {
		ingredient_mapping := schemas.Ingredient_Mapping_Schema{}
		response.Ingredient = &schemas.Intake_Ingredient{}
		// Getting ingredient data
		row := query_ingredient(intake.Ingredient_Mapping_Id, db)
		err = scan_ingredient(row, &ingredient_mapping)
		if err != nil {
			log.Println("Post_Intake | Error on scanning ingredient: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		response.Ingredient.Details = ingredient_mapping
		images, err := get_ingredient_images(db, intake.Ingredient_Mapping_Id)
		if err != nil {
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		response.Ingredient.Images = images
		response.Food = nil
	}
	if intake.Food_Id != 0 {
		food_mapping := schemas.Food_Mapping_Schema{}
		response.Food = &schemas.Intake_Food{}
		// Getting ingredient data
		row := query_food_and_nutrient(intake.Food_Id, db)
		err = scan_food_and_nutrient(row, &food_mapping.Food, &food_mapping.Nutrient)
		if err != nil {
			log.Println("Post_Intake | Error on scanning food: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		response.Food.Details = food_mapping
		images, err := get_food_images(db, intake.Food_Id)
		if err != nil {
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		response.Food.Images = images
		response.Ingredient = nil
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func query_intake(db *sql.DB, user_id uuid.UUID, intake_id uint) *sql.Row {
	row := db.QueryRow(`SELECT
			intake.id,
			intake.account_id,
			intake.date_created,
			COALESCE(intake.ingredient_mapping_id, 0) as ingredient_mapping_id,
			COALESCE(intake.food_id, 0) as food_id,
			intake.amount,
			intake.amount_unit,
			intake.amount_unit_desc,
			intake.serving_size
		FROM intake
		WHERE intake.account_id = $1 AND intake.id = $2`,
		user_id, intake_id,
	)
	return row
}
func scan_intake(row *sql.Row, intake *models.Intake) error {
	err := row.Scan(
		&intake.ID,
		&intake.Account_Id,
		&intake.Date_Created,
		&intake.Ingredient_Mapping_Id,
		&intake.Food_Id,

		&intake.Amount,
		&intake.Amount_Unit,
		&intake.Amount_Unit_Desc,
		&intake.Serving_Size,
	)
	return err
}
func get_ingredient_images(db *sql.DB, ingredient_mapping_id uint) ([]models.Ingredient_Image, error) {
	rows, err := db.Query(`SELECT
			id,
			ingredient_mapping_id,
			name_file,
			name_url,
			amount,
			amount_unit,
			amount_unit_desc
		FROM ingredient_image
		WHERE ingredient_mapping_id = $1`,
		ingredient_mapping_id,
	)
	if err != nil {
		log.Println("error in querying get_ingredient_images: ", err.Error())
		return nil, err
	}
	defer rows.Close()
	ingredient_images := make([]models.Ingredient_Image, 0, 10)
	for rows.Next() {
		var ingredient_img = models.Ingredient_Image{}
		if err := rows.
			Scan(
				&ingredient_img.ID,
				&ingredient_img.Ingredient_Mapping_Id,
				&ingredient_img.Name_File,
				&ingredient_img.Name_URL,
				&ingredient_img.Amount,
				&ingredient_img.Amount_Unit,
				&ingredient_img.Amount_Unit_Desc,
			); err != nil {
			log.Println("Get_Daily_Nutrients_List | error in scanning Daily_Nutrients: ", err.Error())
			return nil, err
		}
		ingredient_images = append(ingredient_images, ingredient_img)
	}
	return ingredient_images, nil
}
func get_food_images(db *sql.DB, food_id uint) ([]models.Food_Image, error) {
	rows, err := db.Query(`SELECT
			id,
			food_id,
			name_file,
			name_url,
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
				&new_image.Name_URL,
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
