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

func Get_Macros_List(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Macros_List | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Get_Macros_List)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Macros_List | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// querying macros
	macros_list, err := query_and_scan_macros_list(db, id, reqData)
	if err != nil && err != sql.ErrNoRows {
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(macros_list)
}

func query_and_scan_macros_list(db *sql.DB, user_id uuid.UUID, reqData *schemas.Req_Get_Macros_List) ([]models.Macros, error) {
	rows, err := db.Query(`SELECT
			id, date_created, calories, protein, carbs, fats, 
			max_calories, max_protein, max_carbs, max_fats,
			activity_lvl_id, diet_plan_id
		FROM macros WHERE account_id = $1
		AND date_created BETWEEN $2 AND $3
		ORDER BY date_Created desc`,
		user_id, reqData.Start_Date.Format(constants.YYYY_MM_DD), reqData.End_Date.Format(constants.YYYY_MM_DD),
	)
	if err != nil {
		log.Println("Get_Macros_List | error in querying macros: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	days := int(math.Floor(reqData.End_Date.Sub(reqData.Start_Date).Hours() / 24))
	macros := make([]models.Macros, 0, days)
	for rows.Next() {
		var new_macros = models.Macros{Account_Id: user_id}
		if err := rows.
			Scan(
				&new_macros.ID,
				&new_macros.Date_Created,
				&new_macros.Calories,
				&new_macros.Protein,
				&new_macros.Carbs,
				&new_macros.Fats,

				&new_macros.Max_Calories,
				&new_macros.Max_Protein,
				&new_macros.Max_Carbs,
				&new_macros.Max_Fats,
				&new_macros.Activity_Lvl_Id,
				&new_macros.Diet_Plan_Id,
			); err != nil {
			log.Println("Get_Macros_List | error in scanning macros: ", err.Error())
			return nil, err
		}
		macros = append(macros, new_macros)
	}
	return macros, nil
}
