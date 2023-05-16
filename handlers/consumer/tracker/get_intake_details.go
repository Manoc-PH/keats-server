package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/tracker"
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
	response := schemas.Res_Get_Intake_Details{}
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
	// TODO add query for food
	if intake.Ingredient_Mapping_Id != 0 {
		ingredient_mapping := schemas.Ingredient_Mapping_Schema{}
		// Getting ingredient data
		row := query_ingredient(intake.Ingredient_Mapping_Id, db)
		err = scan_ingredient(row, &ingredient_mapping)
		if err != nil {
			log.Println("Post_Intake | Error on scanning ingredient: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
		response.Ingredient.Details = ingredient_mapping
		// TODO add query for ingredient images
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
