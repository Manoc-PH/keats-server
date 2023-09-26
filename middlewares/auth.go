package middlewares

import (
	"database/sql"
	"log"
	"server/constants"
	"server/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func AuthMiddleware(c *fiber.Ctx) (*jwt.Token, uuid.UUID, error) {
	cookie := c.Cookies("jwt")

	// parsing token
	token, err := jwt.ParseWithClaims(
		cookie,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(constants.SecretKey), nil
		},
	)
	if err != nil {
		log.Println(err)
		return nil, uuid.Nil, err
	}
	claims := token.Claims.(*jwt.RegisteredClaims)
	owner_id, err := uuid.Parse(claims.Issuer)
	if err != nil {
		log.Println(err)
		return nil, uuid.Nil, err
	}
	return token, owner_id, nil
}

// TODO REVIEW THIS ENDPOINT
func IsAdmin(owner_id uuid.UUID, db *sql.DB) bool {
	var user = models.Account{}
	// checking if user exists
	row := db.QueryRow(`SELECT id, username FROM account_admin WHERE id = $1`, owner_id)
	// scanning and returning error
	if err := row.Scan(&user.ID, &user.Username); err != nil {
		log.Println("IsAdmin | Error in scanning row: ", err.Error())
		return false
	}
	return true
}
