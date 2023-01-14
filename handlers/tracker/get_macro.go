package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Get_Macro(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Macro | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	macros := models.Macro{}
	// querying macro
	row := query_macro(db, id)
	// scanning and returning error
	err = scan_macro(row, &macros)
	// Server Error
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get_Macro | error in scanning macro: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	// Macro doesnt exist yet
	if err != nil && err == sql.ErrNoRows {
		row = query_account_details(db, id)
		account_vitals := models.Account_Vitals{}
		activity_lvl := models.Activity_Lvl{}
		diet_plan := models.Diet_Plan{}
		scan_account_details(row, &account_vitals, &activity_lvl, &diet_plan)
		calories, err := utilities.Calculate_Calories(
			account_vitals.Sex,
			account_vitals.Weight,
			account_vitals.Height,
			activity_lvl.Bmr_Multiplier,
			diet_plan.Calorie_Percentage,
			account_vitals.Birthday,
		)
		if err != nil {
			log.Println("Get_Macro | error in Calculate_Calories: ", err.Error())
			return utilities.Send_Error(c, "An error occured in calculating your calories", fiber.StatusInternalServerError)
		}
		prtn, crbs, fts := utilities.Calculate_Macros(calories, diet_plan.Protein_Percentage, diet_plan.Carbs_Percentage, diet_plan.Fats_Percentage)
		macros.Total_Calories = float32(calories)
		macros.Total_Protein = prtn
		macros.Total_Carbs = crbs
		macros.Total_Fats = fts
		macros.Date_Created = time.Now()
		// TODO INSERT MACRO TO DATABASE
		err = insert_macros(db, &macros, &account_vitals)
		if err != nil {
			log.Println("Get_Macro | error in insert_macros: ", err.Error())
			return utilities.Send_Error(c, "An error occured in calculating your calories", fiber.StatusInternalServerError)
		}
	}
	return c.Status(fiber.StatusOK).JSON(macros)
}

func query_macro(db *sql.DB, user_id uint) *sql.Row {
	row := db.QueryRow(`SELECT
			id, date_created, calories, protein, carbs, fats, 
			total_calories, total_protein, total_carbs, total_fats,
			activity_lvl_id, diet_plan_id
		FROM macro WHERE account_id = $1 AND date_created::date = date $2`,
		user_id, time.Now(),
		// casting timestamp to date
	)
	return row
}
func scan_macro(row *sql.Row, macros *models.Macro) error {
	err := row.Scan(
		macros.ID,
		macros.Date_Created,
		macros.Calories,
		macros.Protein,
		macros.Carbs,
		macros.Fats,

		macros.Total_Calories,
		macros.Total_Protein,
		macros.Total_Carbs,
		macros.Total_Fats,
		macros.Activity_Lvl_Id,
		macros.Diet_Plan_Id,
	)
	return err
}
func query_account_details(db *sql.DB, user_id uint) *sql.Row {
	row := db.QueryRow(`SELECT
			account_vitals.account_id,
			account_vitals.weight,
			account_vitals.height,
			account_vitals.birthday,
			account_vitals.sex,
			account_vitals.activity_lvl_id,
			account_vitals.diet_plan_id,
			activity_lvl.name,
			activity_lvl.bmr_multipler,
			diet_plan.name,
			diet_plan.calorie_percentage,
			diet_plan.protein_percentage,
			diet_plan.fats_percentage,
			diet_plan.carbs_percentage
		FROM account_vitals WHERE account_vitals.account_id = $1
		JOIN activity_lvl ON account_vitals.activity_lvl_id = activity_lvl.id
		JOIN diet_plan ON account_vitals.diet_plan_id = diet_plan.id
		`,
		user_id,
	)
	return row
}
func scan_account_details(
	row *sql.Row,
	account_vitals *models.Account_Vitals,
	activity_lvl *models.Activity_Lvl,
	diet_plan *models.Diet_Plan,
) error {
	err := row.Scan(
		account_vitals.Account_Id,
		account_vitals.Weight,
		account_vitals.Height,
		account_vitals.Birthday,
		account_vitals.Sex,
		account_vitals.Activity_Lvl_Id,
		account_vitals.Diet_Plan_Id,

		activity_lvl.Name,
		activity_lvl.Bmr_Multiplier,

		diet_plan.Name,
		diet_plan.Calorie_Percentage,
		diet_plan.Protein_Percentage,
		diet_plan.Fats_Percentage,
		diet_plan.Carbs_Percentage,
	)
	return err
}
func insert_macros(db *sql.DB, macros *models.Macro, account_vitals *models.Account_Vitals) error {
	row := db.
		QueryRow(`INSERT INTO macro (
			account_id,
			date_created,
			calories,
			protein,
			carbs,
			fats,
			total_calories,
			total_protein,
			total_carbs,
			total_fats,
			activity_lvl_id,
			diet_plan_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
			macros.Account_Id, macros.Date_Created, 0, 0, 0, 0,
			macros.Total_Calories,
			macros.Total_Protein,
			macros.Total_Carbs,
			macros.Total_Fats,
			account_vitals.Activity_Lvl_Id,
			account_vitals.Diet_Plan_Id,
		)
	err := row.Scan(macros.ID)
	if err != nil {
		return err
	}
	return nil
}
