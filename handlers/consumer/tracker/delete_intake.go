package handlers

import (
	"database/sql"
	"log"
	"server/constants"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/tracker"
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
	nutrient := models.Nutrient{}
	daily_nutrients := models.Daily_Nutrients{}
	// TODO OPTIMIZATION: USE GO ROUTINES
	// Querying intake
	row := query_intake(db, owner_id, reqData.Intake_ID)
	err = scan_intake(row, &intake)
	if err == sql.ErrNoRows {
		return utilities.Send_Error(c, "intake not found", fiber.StatusBadRequest)
	}
	if err != nil {
		log.Println("Delete_Intake | Error on scanning food: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	// Not allowing user to delete past intake
	is_intake_today := check_if_date_is_today(intake.Date_Created, time.Now())
	if !is_intake_today {
		log.Println("Delete_Intake | Error: User trying to delete old intake")
		return utilities.Send_Error(c, "cannot delete intake from more than a day ago", fiber.StatusBadRequest)
	}
	// Querying daily intake
	row = query_daily_nutrients(db, owner_id)
	err = scan_daily_nutrients(row, &daily_nutrients)
	if err != nil {
		log.Println("Delete_Intake | Error on scanning daily_nutrients: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	// TODO ADD HANDLER FOR FOOD
	// Querying nutrient of intake
	if intake.Ingredient_Mapping_Id != constants.Empty_UUID {
		row = query_ingredient_nutrient(intake.Ingredient_Mapping_Id, db)
		err = scan_nutrient(row, &nutrient)
		if err != nil {
			log.Println("Delete_Intake | Error on scanning ingredient nutrients: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
	}
	if intake.Food_Id != constants.Empty_UUID {
		row = query_food_nutrient(intake.Food_Id, db)
		err = scan_nutrient(row, &nutrient)
		if err != nil {
			log.Println("Delete_Intake | Error on scanning food nutrients: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
	}
	// Calculate nutrients to be deleted
	nutrients_to_delete := models.Nutrient{}
	calc_nutrients(&nutrients_to_delete, &nutrient, intake.Amount)
	// Delete intake and nutrients
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	err = delete_intake(txn, &intake)
	if err != nil {
		log.Println("Delete_Intake | Error on delete_intake: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	daily_nutrients_to_delete := models.Daily_Nutrients{ID: daily_nutrients.ID}
	scan_and_inverse_daily_nutrients(&daily_nutrients_to_delete, &nutrients_to_delete)
	err = update_daily_nutrients(txn, &daily_nutrients_to_delete)
	if err != nil {
		log.Println("Delete_Intake | Error on delete_intake: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("Delete_Intake (commit) | Error: ", err.Error())
		return err
	}
	// response_data.Deleted_Coins_And_XP = schemas.Added_Coins_And_XP{Coins: coins * -1, XP: xp * -1}
	inverse_nutrient(&nutrients_to_delete)
	response_data.Deleted_Daily_Nutrients = nutrients_to_delete
	response_data.Intake = intake

	return c.Status(fiber.StatusOK).JSON(response_data)
}

func delete_intake(txn *sql.Tx, intake *models.Intake) error {
	_, err := txn.Exec(`DELETE FROM intake WHERE id = $1`, intake.ID)
	if err != nil {
		log.Println("Delete_Intake (delete_intake) | Error: ", err.Error())
		return err
	}
	return nil
}
func scan_and_inverse_daily_nutrients(daily_nutrients *models.Daily_Nutrients, nutrients_to_delete *models.Nutrient) {
	daily_nutrients.Calories = nutrients_to_delete.Calories * -1
	daily_nutrients.Protein = nutrients_to_delete.Protein * -1
	daily_nutrients.Carbs = nutrients_to_delete.Carbs * -1
	daily_nutrients.Fats = nutrients_to_delete.Fats * -1
	daily_nutrients.Trans_Fat = nutrients_to_delete.Trans_Fat * -1
	daily_nutrients.Saturated_Fat = nutrients_to_delete.Saturated_Fat * -1
	daily_nutrients.Sugars = nutrients_to_delete.Sugars * -1
	daily_nutrients.Fiber = nutrients_to_delete.Fiber * -1
	daily_nutrients.Sodium = nutrients_to_delete.Sodium * -1
	daily_nutrients.Iron = nutrients_to_delete.Iron * -1
	daily_nutrients.Calcium = nutrients_to_delete.Calcium * -1
}
func inverse_nutrient(nutrient *models.Nutrient) {
	nutrient.Calories = nutrient.Calories * -1
	nutrient.Protein = nutrient.Protein * -1
	nutrient.Carbs = nutrient.Carbs * -1
	nutrient.Fats = nutrient.Fats * -1
	nutrient.Trans_Fat = nutrient.Trans_Fat * -1
	nutrient.Saturated_Fat = nutrient.Saturated_Fat * -1
	nutrient.Sugars = nutrient.Sugars * -1
	nutrient.Fiber = nutrient.Fiber * -1
	nutrient.Sodium = nutrient.Sodium * -1
	nutrient.Iron = nutrient.Iron * -1
	nutrient.Calcium = nutrient.Calcium * -1
}

// func delete_intake_d_nutrients_and_gamestat(db *sql.DB, d_nutrients_to_delete *models.Nutrient, coins int, xp int, intake *models.Intake) error {
// 	txn, err := db.Begin()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	_, err = txn.Exec(
// 		`UPDATE daily_nutrients SET
// 			calories = calories - $1,
// 			protein = protein - $2,
// 			carbs = carbs - $3,
// 			fats = fats - $4
// 		WHERE id = $5`,
// 		d_nutrients_to_delete.Calories,
// 		d_nutrients_to_delete.Protein,
// 		d_nutrients_to_delete.Carbs,
// 		d_nutrients_to_delete.Fats,
// 		d_nutrients_to_delete.ID,
// 	)
// 	if err != nil {
// 		log.Println("delete_intake_d_nutrients_and_gamestat (update daily_nutrients) | Error: ", err.Error())
// 		return err
// 	}
// 	_, err = txn.Exec(
// 		`UPDATE account_game_stat SET coins = coins - $1, xp = xp - $2 WHERE account_id = $3`,
// 		coins, xp, intake.Account_Id,
// 	)
// 	if err != nil {
// 		log.Println("delete_intake_d_nutrients_and_gamestat (update account_game_stat)| Error: ", err.Error())
// 		return err
// 	}
// 	_, err = txn.Exec(`DELETE FROM intake WHERE id = $1`, intake.ID)
// 	if err != nil {
// 		log.Println("delete_intake_d_nutrients_and_gamestat (update intake)| Error: ", err.Error())
// 		return err
// 	}
// 	err = txn.Commit()
// 	if err != nil {
// 		txn.Rollback()
// 		log.Println("delete_intake_d_nutrients_and_gamestat (commit) | Error: ", err.Error())
// 		return err
// 	}
// 	return nil
// }
