package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Get_Macros(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, Owner_Id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Macros | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	macros := models.Macros{Account_Id: Owner_Id}
	// querying macros
	row := query_macros(db, Owner_Id)
	// scanning and returning error
	err = scan_macros(row, &macros)
	// Macros doesnt exist yet
	if err != nil && err == sql.ErrNoRows {
		account_vitals := models.Account_Vitals{}
		activity_lvl := models.Activity_Lvl{}
		diet_plan := models.Diet_Plan{}
		err := query_and_scan_account_details(db, Owner_Id, &account_vitals, &activity_lvl, &diet_plan)
		if err != nil {
			log.Println("Get_Macros | error in query_and_scan_account_details: ", err.Error())
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
			log.Println("Get_Macros | error in Calculate_Calories: ", err.Error())
			return utilities.Send_Error(c, "An error occured in calculating your calories", fiber.StatusInternalServerError)
		}
		prtn, crbs, fts := utilities.Calculate_Macros(calories, diet_plan.Protein_Percentage, diet_plan.Carbs_Percentage, diet_plan.Fats_Percentage)
		macros.Total_Calories = calories
		macros.Total_Protein = prtn
		macros.Total_Carbs = crbs
		macros.Total_Fats = fts
		macros.Date_Created = time.Now()
		macros.Activity_Lvl_Id = account_vitals.Activity_Lvl_Id
		macros.Diet_Plan_Id = account_vitals.Diet_Plan_Id
		err = insert_macros(db, &macros, &account_vitals)
		if err != nil {
			log.Println("Get_Macros | error in insert_macros: ", err.Error())
			return utilities.Send_Error(c, "An error occured in saving your macros", fiber.StatusInternalServerError)
		}
	}
	// Server Error
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get_Macros | error in scanning macros: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(macros)
}

func query_macros(db *sql.DB, user_id uuid.UUID) *sql.Row {
	row := db.QueryRow(`SELECT
			id, date_created, calories, protein, carbs, fats, 
			total_calories, total_protein, total_carbs, total_fats,
			activity_lvl_id, diet_plan_id
		FROM macros WHERE account_id = $1 AND date_created = $2;`,
		user_id, time.Now().Format("2006-01-02"),
		// casting timestamp to date
	)
	return row
}
func scan_macros(row *sql.Row, macros *models.Macros) error {
	err := row.Scan(
		&macros.ID,
		&macros.Date_Created,
		&macros.Calories,
		&macros.Protein,
		&macros.Carbs,
		&macros.Fats,

		&macros.Total_Calories,
		&macros.Total_Protein,
		&macros.Total_Carbs,
		&macros.Total_Fats,
		&macros.Activity_Lvl_Id,
		&macros.Diet_Plan_Id,
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
func insert_macros(db *sql.DB, macros *models.Macros, account_vitals *models.Account_Vitals) error {
	row := db.
		QueryRow(`INSERT INTO macros (
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
			macros.Account_Id, time.Now().Format("2006-01-02"), 0, 0, 0, 0,
			macros.Total_Calories,
			macros.Total_Protein,
			macros.Total_Carbs,
			macros.Total_Fats,
			account_vitals.Activity_Lvl_Id,
			account_vitals.Diet_Plan_Id,
		)
	err := row.Scan(&macros.ID)
	if err != nil {
		return err
	}
	return nil
}
