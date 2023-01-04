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
	if err := middlewares.Body_Validation(reqData, c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// TODO Query Activity level and Diet Plan tables to verify if both ids sent are valid

	// hashing password and formatting reqData
	password, _ := bcrypt.GenerateFromPassword([]byte(reqData.Password), 10)
	user := models.User{
		Username:        reqData.Username,
		Password:        password,
		Weight:          reqData.Weight,
		Height:          reqData.Height,
		Age:             reqData.Age,
		Sex:             reqData.Sex,
		Activity_Lvl_Id: reqData.Activity_Lvl_Id,
		Diet_Plan_Id:    reqData.Diet_Plan_Id,
	}

	// saving user
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = txn.Exec(
		`INSERT INTO users 
		( id,
			username,
			password,
			updated, 
			created,
			profile_image_link,
			profile_title,
			weight,
			height,
			age,
			sex,
			activity_Lvl_Id,
			diet_plan_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		uuid.New(),
		user.Username,
		user.Password,
		time.Now().Format("YYYY-MM-DD"),
		time.Now().Format("YYYY-MM-DD"),
		nil,
		nil,
		user.Weight,
		user.Height,
		user.Age,
		user.Sex,
		user.Activity_Lvl_Id,
		user.Diet_Plan_Id,
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
	// generating jwt token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    uuid.UUID.String(user.ID),
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

	log.Println("Successfully registered user")
	return c.JSON(user)
}
