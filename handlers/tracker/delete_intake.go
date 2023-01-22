package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/tracker"
	"server/utilities"

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
	macros_curr := models.Macros{}
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
	row = query_macros(db, owner_id)
	err = scan_macros(row, &macros_curr)
	if err != nil {
		log.Println("Delete_Intake | Error on scanning macros: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	macros_to_delete := models.Macros{ID: macros_curr.ID, Account_Id: owner_id}
	calc_macros(&macros_to_delete, &food_nutrient, intake.Amount)
	coins, xp, deductions := utilities.Calc_CnXP_On_Intake(float32(macros_to_delete.Calories), float32(macros_curr.Calories-macros_to_delete.Calories), float32(macros_curr.Max_Calories))
	coins = coins - deductions
	xp = xp - deductions

	err = delete_intake_macro_and_gamestat(db, &macros_to_delete, coins, xp, &intake)
	if err != nil {
		log.Println("Delete_Intake | Error on delete_intake_macro_and_gamestat: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	response_data.Deleted_Coins_And_XP = schemas.Added_Coins_And_XP{Coins: coins, XP: xp}
	response_data.Deleted_Macros = schemas.Deleted_Macros{
		Calories: macros_to_delete.Calories,
		Protein:  macros_to_delete.Protein,
		Carbs:    macros_to_delete.Carbs,
		Fats:     macros_to_delete.Carbs}
	response_data.Intake = intake
	response_data.Food = food

	return c.Status(fiber.StatusOK).JSON(response_data)
}

func delete_intake_macro_and_gamestat(db *sql.DB, macros_to_delete *models.Macros, coins int, xp int, intake *models.Intake) error {
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = txn.Exec(
		`UPDATE macros SET
			calories = calories - $1,
			protein = protein - $2,
			carbs = carbs - $3,
			fats = fats - $4
		WHERE id = $5`,
		macros_to_delete.Calories,
		macros_to_delete.Protein,
		macros_to_delete.Carbs,
		macros_to_delete.Fats,
		macros_to_delete.ID,
	)
	if err != nil {
		log.Println("delete_intake_macro_and_gamestat (update macros) | Error: ", err.Error())
		return err
	}
	_, err = txn.Exec(
		`UPDATE account_game_stat SET coins = coins - $1, xp = xp - $2 WHERE account_id = $3`,
		coins, xp, macros_to_delete.Account_Id,
	)
	if err != nil {
		log.Println("delete_intake_macro_and_gamestat (update account_game_stat)| Error: ", err.Error())
		return err
	}
	_, err = txn.Exec(`DELETE FROM intake WHERE id = $1`, intake.ID)
	if err != nil {
		log.Println("delete_intake_macro_and_gamestat (update intake)| Error: ", err.Error())
		return err
	}
	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("delete_intake_macro_and_gamestat (commit) | Error: ", err.Error())
		return err
	}
	return nil
}
