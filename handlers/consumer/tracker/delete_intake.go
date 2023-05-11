package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/tracker"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Delete_Intake(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Delete_Intake | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Delete_Intake)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Delete_Intake | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// Final response
	response_data := schemas.Res_Delete_Intake{}

	//* data processing
	intake := models.Intake{}
	food := models.Food{}
	food_nutrient := models.Food_Nutrient{}
	d_nutrients_curr := models.Daily_Nutrients{}
	// TODO OPTIMIZATION: USE GO ROUTINES
	row := query_intake_food(reqData.Intake_ID, db)
	err = scan_intake_food(row, &intake, &food, &food_nutrient)
	if err == sql.ErrNoRows {
		return utilities.Send_Error(c, "intake not found", fiber.StatusBadRequest)
	}
	if err != nil {
		log.Println("Delete_Intake | Error on scanning food: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	is_intake_today := check_if_date_is_today(intake.Date_Created, time.Now())
	if !is_intake_today {
		log.Println("Delete_Intake | Error: User trying to delete old intake")
		return utilities.Send_Error(c, "cannot delete intake from more than a day ago", fiber.StatusBadRequest)
	}
	row = query_d_nutrients(db, owner_id)
	err = scan_d_nutrients(row, &d_nutrients_curr)
	if err != nil {
		log.Println("Delete_Intake | Error on scanning daily_nutrients: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	d_nutrients_to_delete := models.Daily_Nutrients{ID: d_nutrients_curr.ID, Account_Id: owner_id}
	calc_d_nutrients(&d_nutrients_to_delete, &food_nutrient, intake.Amount)
	coins, xp, deductions := utilities.Calc_CnXP_On_Intake(float32(d_nutrients_to_delete.Calories), float32(d_nutrients_curr.Calories-d_nutrients_to_delete.Calories), float32(d_nutrients_curr.Max_Calories))
	coins = coins - deductions
	xp = xp - deductions

	err = delete_intake_d_nutrients_and_gamestat(db, &d_nutrients_to_delete, coins, xp, &intake)
	if err != nil {
		log.Println("Delete_Intake | Error on delete_intake_d_nutrients_and_gamestat: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	response_data.Deleted_Coins_And_XP = schemas.Added_Coins_And_XP{Coins: coins * -1, XP: xp * -1}
	response_data.Deleted_Daily_Nutrients = schemas.Deleted_Daily_Nutrients{
		Calories: d_nutrients_to_delete.Calories * -1,
		Protein:  d_nutrients_to_delete.Protein * -1,
		Carbs:    d_nutrients_to_delete.Carbs * -1,
		Fats:     d_nutrients_to_delete.Carbs * -1}
	response_data.Intake = intake
	response_data.Food = food

	return c.Status(fiber.StatusOK).JSON(response_data)
}

func delete_intake_d_nutrients_and_gamestat(db *sql.DB, d_nutrients_to_delete *models.Daily_Nutrients, coins int, xp int, intake *models.Intake) error {
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = txn.Exec(
		`UPDATE daily_nutrients SET
			calories = calories - $1,
			protein = protein - $2,
			carbs = carbs - $3,
			fats = fats - $4
		WHERE id = $5`,
		d_nutrients_to_delete.Calories,
		d_nutrients_to_delete.Protein,
		d_nutrients_to_delete.Carbs,
		d_nutrients_to_delete.Fats,
		d_nutrients_to_delete.ID,
	)
	if err != nil {
		log.Println("delete_intake_d_nutrients_and_gamestat (update daily_nutrients) | Error: ", err.Error())
		return err
	}
	_, err = txn.Exec(
		`UPDATE account_game_stat SET coins = coins - $1, xp = xp - $2 WHERE account_id = $3`,
		coins, xp, d_nutrients_to_delete.Account_Id,
	)
	if err != nil {
		log.Println("delete_intake_d_nutrients_and_gamestat (update account_game_stat)| Error: ", err.Error())
		return err
	}
	_, err = txn.Exec(`DELETE FROM intake WHERE id = $1`, intake.ID)
	if err != nil {
		log.Println("delete_intake_d_nutrients_and_gamestat (update intake)| Error: ", err.Error())
		return err
	}
	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("delete_intake_d_nutrients_and_gamestat (commit) | Error: ", err.Error())
		return err
	}
	return nil
}
