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

// TODO Make sure only the owner can view daily nutrients
// Gets the summary of nutrients for a single day
func Get_Daily_Nutrients(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Daily_Nutrients | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	account_exists := check_account_exists(db, owner_id)
	if account_exists == false {
		log.Println("Get_Daily_Nutrients | error account does not exist")
		return utilities.Send_Error(c, "Account does not exist", fiber.StatusUnauthorized)
	}
	daily_nutrients := models.Daily_Nutrients{Account_Id: owner_id}
	// querying Daily_Nutrients
	row := query_daily_nutrients(db, owner_id)
	err = scan_daily_nutrients(row, &daily_nutrients)
	if err != nil && err == sql.ErrNoRows {
		err = generate_daily_nutrients(db, owner_id, &daily_nutrients)
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
func generate_daily_nutrients(db *sql.DB, owner_id uuid.UUID, daily_nutrients *models.Daily_Nutrients) error {
	consumer_vitals := models.Consumer_Vitals{}
	activity_lvl := models.Activity_Lvl{}
	diet_plan := models.Diet_Plan{}
	err := query_and_scan_account_details(db, owner_id, &consumer_vitals, &activity_lvl, &diet_plan)
	if err != nil {
		log.Println("Get_Daily_Nutrients | error in query_and_scan_account_details: ", err.Error())
		return err
	}
	calories, err := utilities.Calculate_Calories(
		consumer_vitals.Sex,
		consumer_vitals.Weight,
		consumer_vitals.Height,
		activity_lvl.Bmr_Multiplier,
		diet_plan.Calorie_Percentage,
		consumer_vitals.Birthday,
	)
	if err != nil {
		log.Println("Get_Daily_Nutrients | error in Calculate_Calories: ", err.Error())
		return err
	}
	prtn, crbs, fts := utilities.Calculate_Daily_Nutrients(calories, diet_plan.Protein_Percentage, diet_plan.Carbs_Percentage, diet_plan.Fats_Percentage)
	daily_nutrients.ID = uuid.New()
	daily_nutrients.Max_Calories = calories
	daily_nutrients.Max_Protein = prtn
	daily_nutrients.Max_Carbs = crbs
	daily_nutrients.Max_Fats = fts
	daily_nutrients.Date_Created = time.Now()
	daily_nutrients.Activity_Lvl_Id = consumer_vitals.Activity_Lvl_Id
	daily_nutrients.Diet_Plan_Id = consumer_vitals.Diet_Plan_Id
	err = insert_d_nutrients(db, daily_nutrients, &consumer_vitals)
	if err != nil {
		log.Println("Get_Daily_Nutrients | error in insert_Daily_Nutrients: ", err.Error())
		return err
	}
	return nil
}

func query_daily_nutrients(db *sql.DB, owner_id uuid.UUID) *sql.Row {
	row := db.QueryRow(`SELECT
			id, date_created, calories, protein, carbs, fats,
			max_calories, max_protein, max_carbs, max_fats,
			activity_lvl_id, diet_plan_id,
			trans_fat, saturated_fat, sugars, fiber, sodium, iron, calcium
		FROM daily_nutrients WHERE account_id = $1 AND date_created = $2;`,
		owner_id, time.Now().Format(constants.YYYY_MM_DD),
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
	db *sql.DB, owner_id uuid.UUID,
	consumer_vitals *models.Consumer_Vitals,
	activity_lvl *models.Activity_Lvl,
	diet_plan *models.Diet_Plan) error {
	row := db.QueryRow(`SELECT
			consumer_vitals.account_id,
			consumer_vitals.weight,
			consumer_vitals.height,
			consumer_vitals.birthday,
			consumer_vitals.sex,
			consumer_vitals.activity_lvl_id,
			consumer_vitals.diet_plan_id,
			activity_lvl.name,
			activity_lvl.bmr_multiplier,
			diet_plan.name,
			diet_plan.calorie_percentage,
			diet_plan.protein_percentage,
			diet_plan.fats_percentage,
			diet_plan.carbs_percentage
		FROM consumer_vitals
		JOIN activity_lvl ON consumer_vitals.activity_lvl_id = activity_lvl.id
		JOIN diet_plan ON consumer_vitals.diet_plan_id = diet_plan.id
		WHERE consumer_vitals.account_id = $1`,
		owner_id,
	)
	err := row.Scan(
		&consumer_vitals.Account_Id,
		&consumer_vitals.Weight,
		&consumer_vitals.Height,
		&consumer_vitals.Birthday,
		&consumer_vitals.Sex,
		&consumer_vitals.Activity_Lvl_Id,
		&consumer_vitals.Diet_Plan_Id,

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
func insert_d_nutrients(db *sql.DB, daily_nutrients *models.Daily_Nutrients, consumer_vitals *models.Consumer_Vitals) error {
	_, err := db.
		Exec(`INSERT INTO daily_nutrients (
			id,
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
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`,
			daily_nutrients.ID, daily_nutrients.Account_Id, time.Now().Format(constants.YYYY_MM_DD), 0, 0, 0, 0,
			daily_nutrients.Max_Calories,
			daily_nutrients.Max_Protein,
			daily_nutrients.Max_Carbs,
			daily_nutrients.Max_Fats,
			consumer_vitals.Activity_Lvl_Id,
			consumer_vitals.Diet_Plan_Id,
		)
	if err != nil {
		return err
	}
	return nil
}
