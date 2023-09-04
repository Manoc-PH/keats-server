package handlers

import (
	"database/sql"
	"log"
	"server/constants"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/auth"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx, db *sql.DB) error {
	// data validation
	reqData := new(schemas.Req_Login)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Login | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	var user = models.Account_Admin{}

	// checking if user exists
	row := db.QueryRow(`SELECT id, username, password FROM account_admin WHERE username = $1`, reqData.Username)
	// scanning and returning error
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		log.Println("Login | Error in scanning row: ", err.Error())
		return utilities.Send_Error(c, "user does not exist", fiber.StatusBadRequest)
	}

	// checking if password matches user
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(reqData.Password)); err != nil {
		log.Println("Login | Error in comparing password: ", err.Error())
		return utilities.Send_Error(c, "incorrect password", fiber.StatusBadRequest)
	}

	// generating jwt token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    uuid.UUID.String(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)),
	})
	token, err := claims.SignedString([]byte(constants.SecretKey))
	if err != nil {
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
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

	user.Password = nil
	log.Println("Successfully logged user in")
	return c.JSON(user)
}
