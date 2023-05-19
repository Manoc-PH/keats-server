package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/account"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Update_Account_Profile(c *fiber.Ctx, db *sql.DB) error {
	// data validation
	reqData := new(schemas.Req_Update_Account_Profile)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Update_Profile | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// TODO Query Activity level and Diet Plan tables to verify if both ids sent are valid

	// saving user
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = txn.Exec(
		`UPDATE users
			updated= $1, 
			profile_image_link= $2,
			profile_title= $3,
			weight= $4,
			height= $5,
			age= $6,
			sex= $7,
			activity_Lvl_Id= $8,
			diet_plan_id= $9
		WHERE id = $10`,
		time.Now().Format("YYYY-MM-DD"),
		reqData.Profile_Image_Link,
		reqData.Profile_Title,
		reqData.Weight,
		reqData.Height,
		reqData.Age,
		reqData.Sex,
		reqData.Activity_Lvl_Id,
		reqData.Diet_Plan_Id,
		reqData.ID,
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
