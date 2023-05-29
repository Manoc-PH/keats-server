package handlers

import (
	"database/sql"
	"log"
	constants "server/constants/formats"
	"server/middlewares"
	"server/models"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Gets the summary of nutrients for a single day
func Get_Daily_Nutrients(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, Owner_Id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Daily_Nutrients | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	account_exists := check_account_exists(db, Owner_Id)
	if account_exists == false {
		log.Println("Get_Daily_Nutrients | error account does not exist")
		return utilities.Send_Error(c, "Account does not exist", fiber.StatusUnauthorized)
	}
	daily_nutrients := models.Daily_Nutrients{Account_Id: Owner_Id}
	// querying Daily_Nutrients
	row := query_daily_nutrients(db, Owner_Id)
	err = scan_daily_nutrients(row, &daily_nutrients)
	if err != nil && err == sql.ErrNoRows {
		err = generate_daily_nutrients(db, Owner_Id, &daily_nutrients)
		if err != nil {
			log.Println("Get_Daily_Nutrients | error in generate_daily_nutrients: ", err.Error())
			return utilities.Send_Error(c, "An error occured in getting daily nutrients", fiber.StatusInternalServerError)
		}
	}
	// Server Error
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get_Daily_Nutrients | error in scanning Daily_Nutrients: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(daily_nutrients)
}
func check_account_exists(db *sql.DB, Owner_Id uuid.UUID) bool {
	row := db.QueryRow(`SELECT coalesce(id, null) FROM account WHERE id = $1;`,
		Owner_Id,
	)
	exiting_account := models.Account{}
	err := row.Scan(&exiting_account.ID)
	if err != nil {
		log.Println("check_account_exists | error in scanning exiting_account: ", err.Error())
		return false
	}
	return true
}
func generate_daily_nutrients(db *sql.DB, Owner_Id uuid.UUID, daily_nutrients *models.Daily_Nutrients) error {
	account_vitals := models.Account_Vitals{}
	activity_lvl := models.Activity_Lvl{}
	diet_plan := models.Diet_Plan{}
	err := query_and_scan_account_details(db, Owner_Id, &account_vitals, &activity_lvl, &diet_plan)
	if err != nil {
		log.Println("Get_Daily_Nutrients | error in query_and_scan_account_details: ", err.Error())
		return err
	}
	calories, err := utilities.Calculate_Calories(
		account_vitals.Sex,
		account_vitals.Weight,
		account_vitals.Height,
		activity_lvl.Bmr_Multiplier,
		diet_plan.Calorie_Percentage,
		account_vitals.Birthday,
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
	daily_nutrients.Date_Created = time.Now()
	daily_nutrients.Activity_Lvl_Id = account_vitals.Activity_Lvl_Id
	daily_nutrients.Diet_Plan_Id = account_vitals.Diet_Plan_Id
	err = insert_d_nutrients(db, daily_nutrients, &account_vitals)
	if err != nil {
		log.Println("Get_Daily_Nutrients | error in insert_Daily_Nutrients: ", err.Error())
		return err
	}
	return nil
}

func query_daily_nutrients(db *sql.DB, user_id uuid.UUID) *sql.Row {
	row := db.QueryRow(`SELECT
			id, date_created, calories, protein, carbs, fats,
			max_calories, max_protein, max_carbs, max_fats,
			activity_lvl_id, diet_plan_id,
			trans_fat, saturated_fat, sugars, fiber, sodium, iron, calcium
		FROM daily_nutrients WHERE account_id = $1 AND date_created = $2;`,
		user_id, time.Now().Format(constants.YYYY_MM_DD),
	// casting timestamp to date
	)
	return row
}
func scan_daily_nutrients(row *sql.Row, daily_nutrients *models.Daily_Nutrients) error {
	err := row.Scan(
		&daily_nutrients.ID,
		&daily_nutrients.Date_Created,
		&daily_nutrients.Calories,
		&daily_nutrients.Protein,
		&daily_nutrients.Carbs,
		&daily_nutrients.Fats,

		&daily_nutrients.Max_Calories,
		&daily_nutrients.Max_Protein,
		&daily_nutrients.Max_Carbs,
		&daily_nutrients.Max_Fats,
		&daily_nutrients.Activity_Lvl_Id,
		&daily_nutrients.Diet_Plan_Id,

		&daily_nutrients.Trans_Fat,
		&daily_nutrients.Saturated_Fat,
		&daily_nutrients.Sugars,
		&daily_nutrients.Fiber,
		&daily_nutrients.Sodium,
		&daily_nutrients.Iron,
		&daily_nutrients.Calcium,
	)
	return err
}
func query_and_scan_account_details(
	db *sql.DB, user_id uuid.UUID,
	account_vitals *models.Account_Vitals,
	activity_lvl *models.Activity_Lvl,
	diet_plan *models.Diet_Plan) error {
	row := db.QueryRow(`SELECT
			account_vitals.account_id,
			account_vitals.weight,
			account_vitals.height,
			account_vitals.birthday,
			account_vitals.sex,
			account_vitals.activity_lvl_id,
			account_vitals.diet_plan_id,
			activity_lvl.name,
			activity_lvl.bmr_multiplier,
			diet_plan.name,
			diet_plan.calorie_percentage,
			diet_plan.protein_percentage,
			diet_plan.fats_percentage,
			diet_plan.carbs_percentage
		FROM account_vitals
		JOIN activity_lvl ON account_vitals.activity_lvl_id = activity_lvl.id
		JOIN diet_plan ON account_vitals.diet_plan_id = diet_plan.id
		WHERE account_vitals.account_id = $1`,
		user_id,
	)
	err := row.Scan(
		&account_vitals.Account_Id,
		&account_vitals.Weight,
		&account_vitals.Height,
		&account_vitals.Birthday,
		&account_vitals.Sex,
		&account_vitals.Activity_Lvl_Id,
		&account_vitals.Diet_Plan_Id,

		&activity_lvl.Name,
		&activity_lvl.Bmr_Multiplier,

		&diet_plan.Name,
		&diet_plan.Calorie_Percentage,
		&diet_plan.Protein_Percentage,
		&diet_plan.Fats_Percentage,
		&diet_plan.Carbs_Percentage,
	)
	return err
}
func insert_d_nutrients(db *sql.DB, daily_nutrients *models.Daily_Nutrients, account_vitals *models.Account_Vitals) error {
	row := db.
		QueryRow(`INSERT INTO daily_nutrients (
			account_id,
			date_created,
			calories,
			protein,
			carbs,
			fats,
			max_calories,
			max_protein,
			max_carbs,
			max_fats,
			activity_lvl_id,
			diet_plan_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
			daily_nutrients.Account_Id, time.Now().Format(constants.YYYY_MM_DD), 0, 0, 0, 0,
			daily_nutrients.Max_Calories,
			daily_nutrients.Max_Protein,
			daily_nutrients.Max_Carbs,
			daily_nutrients.Max_Fats,
			account_vitals.Activity_Lvl_Id,
			account_vitals.Diet_Plan_Id,
		)
	err := row.Scan(&daily_nutrients.ID)
	if err != nil {
		return err
	}
	return nil
}
