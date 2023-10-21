package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Patch_Recipe(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Recipe | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Patch_Recipe)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Recipe | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}
