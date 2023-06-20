package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/common"

	"github.com/gofiber/fiber/v2"
)

func Get_Name_Availability(c *fiber.Ctx, db *sql.DB) error {
	//* data validation
	reqData := new(schemas.Req_Get_Name_Availability)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Name_Availability | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	res := fetch_name(db, reqData.Username)

	return c.JSON(res)
}

func fetch_name(db *sql.DB, username string) bool {
	user := models.Account{}
	taken := true
	row := db.QueryRow(`SELECT username FROM account WHERE username = $1`, username)
	err := row.Scan(&user.Username)
	if err != nil {
		log.Println("Get_Name_Availability | error in querying Account: ", err.Error())
		return false
	}
	if user.Username != "" {
		taken = true
	}
	return taken
}
