package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/account"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Update_Consumer_Vitals(c *fiber.Ctx, db *sql.DB) error {
	// data validation
	reqData := new(schemas.Req_Update_Consumer_Vitals)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Update_Vitals | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// TODO Query Activity level and Diet Plan tables to verify if both ids sent are valid
	Found_Activity_lvl := models.Activity_Lvl{}
	Found_Diet_Plan := models.Diet_Plan{}
	err := query_activity_level(db, reqData.Activity_Lvl_Id, &Found_Activity_lvl)
	if err != nil {
		log.Println("Update_Vitals | Error on query validation: ", err.Error())
		return utilities.Send_Error(c, "Activity Level does not exist", fiber.StatusBadRequest)
	}
	err = query_diet_plan(db, reqData.Diet_Plan_Id, &Found_Diet_Plan)
	if err != nil {
		log.Println("Update_Vitals | Error on query validation: ", err.Error())
		return utilities.Send_Error(c, "Diet Plan does not exist", fiber.StatusBadRequest)
	}
	// saving user
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = txn.Exec(
		`UPDATE consumer_vitals SET
			weight = $1,
			height = $2,
			birthday = $3,
			sex = $4,
			activity_lvl_id = $5,
			diet_plan_id = $6
		WHERE account_id = $7`,
		reqData.Weight,
		reqData.Height,
		reqData.Birthday,
		reqData.Sex,
		reqData.Activity_Lvl_Id,
		reqData.Diet_Plan_Id,
		reqData.Account_ID,
	)
	if err != nil {
		log.Println("Error updating consumer vitals: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	daily_nutrients := models.Daily_Nutrients{}
	err = update_daily_nutrients(txn, reqData, &Found_Activity_lvl, &Found_Diet_Plan, &daily_nutrients)
	if err != nil {
		log.Println("Error updating daily nutrients: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("Error committing: ", err.Error())
		return err
	}

	log.Println("Successfully updated user's vitals")
	res := schemas.Res_Update_Consumer_Vitals{ReqData: *reqData, Daily_Nutrients: daily_nutrients}
	return c.JSON(res)
}

func update_daily_nutrients(
	txn *sql.Tx,
	consumer_vitals *schemas.Req_Update_Consumer_Vitals,
	activity_lvl *models.Activity_Lvl,
	diet_plan *models.Diet_Plan,
	daily_nutrients *models.Daily_Nutrients,
) error {
	calories, err := utilities.Calculate_Calories(
		consumer_vitals.Sex,
		int(consumer_vitals.Weight),
		int(consumer_vitals.Height),
		activity_lvl.Bmr_Multiplier,
		diet_plan.Calorie_Percentage,
		consumer_vitals.Birthday,
	)
	if err != nil {
		log.Println("Get_Daily_Nutrients | error in Calculate_Calories: ", err.Error())
		return err
	}

	prtn, crbs, fts := utilities.Calculate_Daily_Nutrients(calories, diet_plan.Protein_Percentage, diet_plan.Carbs_Percentage, diet_plan.Fats_Percentage)
	daily_nutrients.Max_Calories = calories
	daily_nutrients.Max_Protein = prtn
	daily_nutrients.Max_Carbs = crbs
	daily_nutrients.Max_Fats = fts
	daily_nutrients.Activity_Lvl_Id = consumer_vitals.Activity_Lvl_Id
	daily_nutrients.Diet_Plan_Id = consumer_vitals.Diet_Plan_Id

	_, err = txn.Exec(
		`UPDATE daily_nutrients SET
			max_calories = $1,
			max_protein = $2,
			max_carbs = $3,
			max_fats = $4,
			activity_lvl_id = $5,
			diet_plan_id = $6
		WHERE account_id = $7 AND date_created = $8`,
		daily_nutrients.Max_Calories,
		daily_nutrients.Max_Protein,
		daily_nutrients.Max_Carbs,
		daily_nutrients.Max_Fats,
		daily_nutrients.Activity_Lvl_Id,
		daily_nutrients.Diet_Plan_Id,
		consumer_vitals.Account_ID,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return err
	}
	return nil
}

func query_activity_level(db *sql.DB, activity_id uuid.UUID, activity_lvl *models.Activity_Lvl) error {
	row := db.QueryRow(`SELECT id, bmr_multiplier FROM activity_lvl WHERE id = $1`, activity_id)
	err := row.Scan(&activity_lvl.ID, &activity_lvl.Bmr_Multiplier)
	return err
}
func query_diet_plan(db *sql.DB, diet_plan_id uuid.UUID, diet_plan *models.Diet_Plan) error {
	row := db.QueryRow(`
		SELECT 
			id, calorie_percentage, protein_percentage, fats_percentage, carbs_percentage
		FROM diet_plan WHERE id = $1`, diet_plan_id)
	err := row.Scan(
		&diet_plan.ID,
		&diet_plan.Calorie_Percentage,
		&diet_plan.Protein_Percentage,
		&diet_plan.Fats_Percentage,
		&diet_plan.Carbs_Percentage)
	return err
}
