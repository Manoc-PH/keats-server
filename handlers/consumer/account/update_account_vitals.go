package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/account"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Update_Account_Vitals(c *fiber.Ctx, db *sql.DB) error {
	// data validation
	reqData := new(schemas.Req_Update_Account_Vitals)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Update_Vitals | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// TODO Query Activity level and Diet Plan tables to verify if both ids sent are valid

	// saving user
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = txn.Exec(
		`UPDATE account_vitals
			weight = $1,
			height = $2,
			birthday = $3,
			sex = $4,
			activity_lvl_id = $5,
			diet_plan_id = $6
		WHERE account_id = $7`,
		time.Now().Format("YYYY-MM-DD"),
		reqData.Weight,
		reqData.Height,
		reqData.Birthday,
		reqData.Sex,
		reqData.Activity_Lvl_Id,
		reqData.Diet_Plan_Id,
		reqData.Account_ID,
	)

	if err != nil {
		log.Println("Error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("Error: ", err.Error())
		return err
	}

	log.Println("Successfully updated user")
	return c.JSON(reqData)
}
