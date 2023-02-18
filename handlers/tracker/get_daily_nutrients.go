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

func Get_Daily_Nutrients(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, Owner_Id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Daily_Nutrients | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	d_nutrients := models.Daily_Nutrients{Account_Id: Owner_Id}
	// querying Daily_Nutrients
	row := query_d_nutrients(db, Owner_Id)
	// scanning and returning error
	err = scan_d_nutrients(row, &d_nutrients)
	// Daily_Nutrients doesnt exist yet
	if err != nil && err == sql.ErrNoRows {
		account_vitals := models.Account_Vitals{}
		activity_lvl := models.Activity_Lvl{}
		diet_plan := models.Diet_Plan{}
		err := query_and_scan_account_details(db, Owner_Id, &account_vitals, &activity_lvl, &diet_plan)
		if err != nil {
			log.Println("Get_Daily_Nutrients | error in query_and_scan_account_details: ", err.Error())
			return utilities.Send_Error(c, "An error occured in getting account details", fiber.StatusInternalServerError)
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
			return utilities.Send_Error(c, "An error occured in calculating your calories", fiber.StatusInternalServerError)
		}
		prtn, crbs, fts := utilities.Calculate_Daily_Nutrients(calories, diet_plan.Protein_Percentage, diet_plan.Carbs_Percentage, diet_plan.Fats_Percentage)
		d_nutrients.Max_Calories = calories
		d_nutrients.Max_Protein = prtn
		d_nutrients.Max_Carbs = crbs
		d_nutrients.Max_Fats = fts
		d_nutrients.Date_Created = time.Now()
		d_nutrients.Activity_Lvl_Id = account_vitals.Activity_Lvl_Id
		d_nutrients.Diet_Plan_Id = account_vitals.Diet_Plan_Id
		err = insert_d_nutrients(db, &d_nutrients, &account_vitals)
		if err != nil {
			log.Println("Get_Daily_Nutrients | error in insert_Daily_Nutrients: ", err.Error())
			return utilities.Send_Error(c, "An error occured in saving your Daily_Nutrients", fiber.StatusInternalServerError)
		}
	}
	// Server Error
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get_Daily_Nutrients | error in scanning Daily_Nutrients: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(d_nutrients)
}

func query_d_nutrients(db *sql.DB, user_id uuid.UUID) *sql.Row {
	row := db.QueryRow(`SELECT
			id, date_created, calories, protein, carbs, fats,
			max_calories, max_protein, max_carbs, max_fats,
			activity_lvl_id, diet_plan_id
		FROM daily_nutrients WHERE account_id = $1 AND date_created = $2;`,
		user_id, time.Now().Format(constants.YYYY_MM_DD),
	// casting timestamp to date
	)
	return row
}
func scan_d_nutrients(row *sql.Row, d_nutrients *models.Daily_Nutrients) error {
	err := row.Scan(
		&d_nutrients.ID,
		&d_nutrients.Date_Created,
		&d_nutrients.Calories,
		&d_nutrients.Protein,
		&d_nutrients.Carbs,
		&d_nutrients.Fats,

		&d_nutrients.Max_Calories,
		&d_nutrients.Max_Protein,
		&d_nutrients.Max_Carbs,
		&d_nutrients.Max_Fats,
		&d_nutrients.Activity_Lvl_Id,
		&d_nutrients.Diet_Plan_Id,
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
