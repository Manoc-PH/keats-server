package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/account"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Get_Account_Vitals(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Intake | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	res := schemas.Res_Get_Account_Vitals{}
	row := query_account_vitals(db, owner_id)
	err = scan_account_vitals(row, &res)
	if err != nil {
		log.Println("Error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.JSON(res)
}

func query_account_vitals(db *sql.DB, user_id uuid.UUID) *sql.Row {
	row := db.QueryRow(`SELECT
			account_vitals.account_id,
			account_vitals.weight,
			account_vitals.height,
			account_vitals.birthday,
			account_vitals.sex,
			account_vitals.activity_lvl_id,
			account_vitals.diet_plan_id,
			activity_lvl.name as activity_lvl_name,
			activity_lvl.bmr_multiplier,
			diet_plan.id as diet_plan_id,
			diet_plan.name as diet_plan_name,
			diet_plan.calorie_percentage,
			diet_plan.protein_percentage,
			diet_plan.fats_percentage,
			diet_plan.carbs_percentage
		FROM account_vitals
		JOIN activity_lvl ON account_vitals.activity_lvl_id = activity_lvl.id
		JOIN diet_plan ON account_vitals.diet_plan_id = diet_plan.id 
		WHERE account_vitals.account_id = $1`, user_id,
	)
	return row
}

func scan_account_vitals(row *sql.Row, account_vitals *schemas.Res_Get_Account_Vitals) error {
	err := row.Scan(
		&account_vitals.Account_ID,
		&account_vitals.Weight,
		&account_vitals.Height,
		&account_vitals.Birthday,
		&account_vitals.Sex,
		&account_vitals.Activity_Lvl_Id,
		&account_vitals.Diet_Plan_Id,
		&account_vitals.Activity_Lvl_Name,
		&account_vitals.Bmr_Multiplier,
		&account_vitals.Diet_Plan_Id,
		&account_vitals.Diet_Plan_Name,
		&account_vitals.Calorie_Percentage,
		&account_vitals.Protein_Percentage,
		&account_vitals.Fats_Percentage,
		&account_vitals.Carbs_Percentage,
	)
	return err
}
