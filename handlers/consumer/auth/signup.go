package handlers

import (
	"database/sql"
	"log"
	"server/constants"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/auth"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Sign_Up(c *fiber.Ctx, db *sql.DB) error {
	// data validation
	reqData := new(schemas.Req_Sign_Up)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Sign_Up | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// TODO Query Activity level and Diet Plan tables to verify if both ids sent are valid

	// hashing password and formatting reqData
	password, _ := bcrypt.GenerateFromPassword([]byte(reqData.Password), 10)
	account_id := uuid.New()
	account_type_id, err := uuid.Parse("4c3c69b0-2eae-4b3c-80e1-619f4718d272")
	if err != nil {
		log.Println("Sign_Up | Error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	account_vitals := models.Account_Vitals{
		ID:              uuid.New(),
		Account_Id:      account_id,
		Weight:          reqData.Weight,
		Height:          reqData.Height,
		Birthday:        reqData.Birthday,
		Sex:             reqData.Sex,
		Activity_Lvl_Id: reqData.Activity_Lvl_Id,
		Diet_Plan_Id:    reqData.Diet_Plan_Id,
	}
	account_profile := models.Consumer_Profile{
		ID:                uuid.New(),
		Account_Id:        account_id,
		Date_Updated:      time.Now(),
		Date_Created:      time.Now(),
		Account_Vitals_Id: account_vitals.ID,
		Measure_Unit_Id:   reqData.Measure_Unit_Id,
	}
	account := models.Account{
		ID:              account_id,
		Username:        reqData.Username,
		Password:        password,
		Account_Type_Id: account_type_id,
	}

	// saving account
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = txn.Exec(
		`INSERT INTO account_vitals 
		( id,
			account_id,
			weight,
			height,
			birthday,
			sex,
			activity_lvl_id,
			diet_plan_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		account_vitals.ID,
		account_vitals.Account_Id,
		account_vitals.Weight,
		account_vitals.Height,
		account_vitals.Birthday,
		account_vitals.Sex,
		account_vitals.Activity_Lvl_Id,
		account_vitals.Diet_Plan_Id,
	)
	if err != nil {
		log.Println("Sign_Up | Error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	_, err = txn.Exec(
		`INSERT INTO consumer_profile (
			id,
			account_id,
			date_updated,
			date_created,
			account_vitals_id,
			measure_unit_id)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		account_profile.ID,
		account_profile.Account_Id,
		account_profile.Date_Updated,
		account_profile.Date_Created,
		account_profile.Account_Vitals_Id,
		account_profile.Measure_Unit_Id,
	)
	if err != nil {
		log.Println("Sign_Up | Error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	_, err = txn.Exec(
		`INSERT INTO account 
		(id, username, password, account_type_id)
		VALUES ($1, $2, $3, $4)`,
		account.ID,
		account.Username,
		account.Password,
		account.Account_Type_Id,
	)
	if err != nil {
		log.Println("Sign_Up | Error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("Sign_Up | Error: ", err.Error())
		return err
	}
	// generating jwt token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    uuid.UUID.String(account.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)),
	})
	token, err := claims.SignedString([]byte(constants.SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	// saving jwt to cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().AddDate(0, 1, 0),
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
	}

	c.Cookie(&cookie)

	log.Println("Successfully registered account")
	return c.JSON(account)
}
