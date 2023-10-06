package handlers

import (
	"database/sql"
	"log"
	"math"
	constants "server/constants/formats"
	"server/middlewares"
	schemas "server/schemas/consumer/tracker"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Gets the summary of daily calories through a date range
func Get_Daily_Calorie_List(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Daily_Calorie_List | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Get_Daily_Calorie_List)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Daily_Calorie_List | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// querying Daily_Nutrients
	d_nutrients_list, err := query_and_scan_d_nutrients_list(db, id, reqData)
	if err != nil && err != sql.ErrNoRows {
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(d_nutrients_list)
}

func query_and_scan_d_nutrients_list(db *sql.DB, user_id uuid.UUID, reqData *schemas.Req_Get_Daily_Calorie_List) ([]schemas.Res_Get_Daily_Calorie_List, error) {
	rows, err := db.Query(`SELECT
			id, date_created, calories
		FROM daily_nutrients WHERE account_id = $1
		AND date_created BETWEEN $2 AND $3
		ORDER BY date_created desc`,
		user_id, reqData.Start_Date.Format(constants.YYYY_MM_DD), reqData.End_Date.Format(constants.YYYY_MM_DD),
	)
	if err != nil {
		log.Println("Get_Daily_Calorie_List | error in querying Daily_Nutrients: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	days := int(math.Floor(reqData.End_Date.Sub(reqData.Start_Date).Hours() / 24))
	daily_nutrients := make([]schemas.Res_Get_Daily_Calorie_List, 0, days)
	for rows.Next() {
		var new_daily_nutrient = schemas.Res_Get_Daily_Calorie_List{Account_Id: user_id}
		if err := rows.
			Scan(
				&new_daily_nutrient.ID,
				&new_daily_nutrient.Date_Created,
				&new_daily_nutrient.Calories,
			); err != nil {
			log.Println("Get_Daily_Calorie_List | error in scanning Daily_Nutrients: ", err.Error())
			return nil, err
		}
		daily_nutrients = append(daily_nutrients, new_daily_nutrient)
	}
	return daily_nutrients, nil
}
