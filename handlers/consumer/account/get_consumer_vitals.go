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

func Get_Consumer_Vitals(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Intake | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	res := schemas.Res_Get_Consumer_Vitals{}
	row := query_consumer_vitals(db, owner_id)
	err = scan_consumer_vitals(row, &res)
	if err != nil {
		log.Println("Error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.JSON(res)
}

func query_consumer_vitals(db *sql.DB, user_id uuid.UUID) *sql.Row {
	row := db.QueryRow(`SELECT
			consumer_vitals.account_id,
			consumer_vitals.weight,
			consumer_vitals.height,
			consumer_vitals.birthday,
			consumer_vitals.sex,
			consumer_vitals.activity_lvl_id,
			consumer_vitals.diet_plan_id,
			activity_lvl.name as activity_lvl_name,
			activity_lvl.bmr_multiplier,
			diet_plan.id as diet_plan_id,
			diet_plan.name as diet_plan_name,
			diet_plan.calorie_percentage,
			diet_plan.protein_percentage,
			diet_plan.fats_percentage,
			diet_plan.carbs_percentage
		FROM consumer_vitals
		JOIN activity_lvl ON consumer_vitals.activity_lvl_id = activity_lvl.id
		JOIN diet_plan ON consumer_vitals.diet_plan_id = diet_plan.id 
		WHERE consumer_vitals.account_id = $1`, user_id,
	)
	return row
}

func scan_consumer_vitals(row *sql.Row, consumer_vitals *schemas.Res_Get_Consumer_Vitals) error {
	err := row.Scan(
		&consumer_vitals.Account_ID,
		&consumer_vitals.Weight,
		&consumer_vitals.Height,
		&consumer_vitals.Birthday,
		&consumer_vitals.Sex,
		&consumer_vitals.Activity_Lvl_Id,
		&consumer_vitals.Diet_Plan_Id,
		&consumer_vitals.Activity_Lvl_Name,
		&consumer_vitals.Bmr_Multiplier,
		&consumer_vitals.Diet_Plan_Id,
		&consumer_vitals.Diet_Plan_Name,
		&consumer_vitals.Calorie_Percentage,
		&consumer_vitals.Protein_Percentage,
		&consumer_vitals.Fats_Percentage,
		&consumer_vitals.Carbs_Percentage,
	)
	return err
}
