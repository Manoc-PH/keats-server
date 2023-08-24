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
	row := db.QueryRow(`
		SELECT account.id, account_type.name
		FROM account
		JOIN account_type ON account.account_type_id = account_type.id
		WHERE account.id = $1;
	`, owner_id)
	account := models.Account{}
	account_type := models.Account_Type{}

	err := row.Scan(&account.ID, &account_type.Name)
	if err != nil {
		return false
	}
	if account.ID != owner_id {
		return false
	}
	if account_type.Name != constants.Account_Types.Admin {
		return false
	}
	return true
}
