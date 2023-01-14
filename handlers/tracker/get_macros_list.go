package handlers

import (
	"database/sql"
	"log"
	"math"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/tracker"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
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
	if err = middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Macros_List | Error on query validation: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusBadRequest)
	}
	days := int(math.Floor(reqData.End_Date.Sub(reqData.Start_Date).Hours() / 24))

	macros := make([]models.Macros, days)
	// querying macros
	rows, err := query_macros_list(db, id, reqData)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get_Macros_List | error in querying macros: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	defer rows.Close()
	// scanning and returning error
	err = scan_macros_list(rows, macros)
	// Server Error
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get_Macros_List | error in scanning macros: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(macros)
}

func query_macros_list(db *sql.DB, user_id uint, reqData *schemas.Req_Get_Macros_List) (*sql.Rows, error) {
	row, err := db.Query(`SELECT
			id, date_created, calories, protein, carbs, fats, 
			total_calories, total_protein, total_carbs, total_fats,
			activity_lvl_id, diet_plan_id
		FROM macros WHERE account_id = $1
		AND date_created::date BETWEEN $2 AND $3
		ORDER BY date_Created asc`,
		// casting timestamp to date
		user_id, reqData.Start_Date, reqData.End_Date,
	)
	return row, err
}
func scan_macros_list(rows *sql.Rows, macros []models.Macros) error {
	for rows.Next() {
		new_macros := models.Macros{}
		if err := rows.
			Scan(
				&new_macros.ID,
				&new_macros.Date_Created,
				&new_macros.Calories,
				&new_macros.Protein,
				&new_macros.Carbs,
				&new_macros.Fats,

				&new_macros.Total_Calories,
				&new_macros.Total_Protein,
				&new_macros.Total_Carbs,
				&new_macros.Total_Fats,
				&new_macros.Activity_Lvl_Id,
				&new_macros.Diet_Plan_Id,
			); err != nil {
			return err
		}
		macros = append(macros, new_macros)
	}
	return nil
}
