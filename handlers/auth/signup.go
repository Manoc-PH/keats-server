package handlers

import (
	"database/sql"
	"kryptoverse-api/constants"
	"kryptoverse-api/middlewares"
	"kryptoverse-api/models"
	schemas "kryptoverse-api/schemas/auth"
	"log"
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

	// hashing password and formatting reqData
	password, _ := bcrypt.GenerateFromPassword([]byte(reqData.Password), 10)
	user := models.User{
		Username: reqData.Username,
		Password: password,
	}

	// saving user
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	row := txn.
		QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", user.Username, user.Password)

	err = row.Scan(&user.ID)
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
