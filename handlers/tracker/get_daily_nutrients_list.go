package handlers

import (
	"database/sql"
	"log"
	"math"
	constants "server/constants/formats"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/tracker"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Gets the summary of daily nutrients through a date range
func Get_Daily_Nutrients_List(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Daily_Nutrients_List | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Get_Daily_Nutrients_List)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Daily_Nutrients_List | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// querying Daily_Nutrients
	d_nutrients_list, err := query_and_scan_d_nutrients_list(db, id, reqData)
	if err != nil && err != sql.ErrNoRows {
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(d_nutrients_list)
}

func query_and_scan_d_nutrients_list(db *sql.DB, user_id uuid.UUID, reqData *schemas.Req_Get_Daily_Nutrients_List) ([]models.Daily_Nutrients, error) {
	rows, err := db.Query(`SELECT
			id, date_created, calories, protein, carbs, fats, 
			max_calories, max_protein, max_carbs, max_fats,
			activity_lvl_id, diet_plan_id
		FROM daily_nutrients WHERE account_id = $1
		AND date_created BETWEEN $2 AND $3
		ORDER BY date_Created desc`,
		user_id, reqData.Start_Date.Format(constants.YYYY_MM_DD), reqData.End_Date.Format(constants.YYYY_MM_DD),
	)
	if err != nil {
		log.Println("Get_Daily_Nutrients_List | error in querying Daily_Nutrients: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	days := int(math.Floor(reqData.End_Date.Sub(reqData.Start_Date).Hours() / 24))
	daily_nutrients := make([]models.Daily_Nutrients, 0, days)
	for rows.Next() {
		var new_d_nutrient = models.Daily_Nutrients{Account_Id: user_id}
		if err := rows.
			Scan(
				&new_d_nutrient.ID,
				&new_d_nutrient.Date_Created,
				&new_d_nutrient.Calories,
				&new_d_nutrient.Protein,
				&new_d_nutrient.Carbs,
				&new_d_nutrient.Fats,

				&new_d_nutrient.Max_Calories,
				&new_d_nutrient.Max_Protein,
				&new_d_nutrient.Max_Carbs,
				&new_d_nutrient.Max_Fats,
				&new_d_nutrient.Activity_Lvl_Id,
				&new_d_nutrient.Diet_Plan_Id,
			); err != nil {
			log.Println("Get_Daily_Nutrients_List | error in scanning Daily_Nutrients: ", err.Error())
			return nil, err
		}
		daily_nutrients = append(daily_nutrients, new_d_nutrient)
	}
	return daily_nutrients, nil
}
