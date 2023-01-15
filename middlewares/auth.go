package middlewares

import (
	"log"
	"server/constants"

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
