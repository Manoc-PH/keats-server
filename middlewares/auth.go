package middlewares

import (
	"kryptoverse-api/constants"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(c *fiber.Ctx) (*jwt.Token, uint, error) {
	SecretKey := os.Getenv("SECRET_KEY")
	cookie := c.Cookies("jwt")

	// parsing token
	token, err := jwt.ParseWithClaims(
		cookie,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}
	claims := token.Claims.(*jwt.RegisteredClaims)
	u64, err := strconv.ParseUint(claims.Issuer, 10, 32)
	if err != nil {
		return nil, 0, err
	}

	owner_id := uint(u64)
	return token, owner_id, nil
}

func Socket_Auth_Middleware(jwt_token *string) (*jwt.Token, uint, error) {
	log.Println(*jwt_token)
	// parsing token
	token, err := jwt.ParseWithClaims(
		*jwt_token,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(constants.SecretKey), nil
		},
	)
	// TODO ADD EXPIRATION HANDLER
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}
	claims := token.Claims.(*jwt.RegisteredClaims)
	u64, err := strconv.ParseUint(claims.Issuer, 10, 32)
	if err != nil {
		return nil, 0, err
	}

	owner_id := uint(u64)
	return token, owner_id, nil
}
