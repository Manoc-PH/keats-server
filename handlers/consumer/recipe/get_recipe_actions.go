package handlers

import (
	"database/sql"
	"log"
	"server/constants"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Get_Recipe_Actions(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Recipe_Discovery | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	} //* data validation
	reqData := new(schemas.Req_Get_Actions)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("get_recipe_details | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	liked := is_recipe_liked(db, reqData.Recipe_Id, owner_id)
	reviewed := is_recipe_reviewed(db, reqData.Recipe_Id, owner_id)

	response := schemas.Res_Get_Actions{Liked: liked, Reviewed: reviewed}
	return c.Status(fiber.StatusOK).JSON(response)
}

func is_recipe_liked(db *sql.DB, recipe_id uuid.UUID, user_id uuid.UUID) bool {
	row := db.QueryRow(`SELECT id FROM recipe_like WHERE owner_id = $1 AND recipe_id = $2 `, user_id, recipe_id)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			log.Println("ERROR on is_recipe_liked: ", err)
		}
		return false
	}
	if id == constants.Empty_UUID {
		return false
	}
	return true
}
func is_recipe_reviewed(db *sql.DB, recipe_id uuid.UUID, user_id uuid.UUID) bool {
	row := db.QueryRow(`SELECT id FROM recipe_review WHERE owner_id = $1 AND recipe_id = $2 `, user_id, recipe_id)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			log.Println("ERROR on is_recipe_reviewed: ", err)
		}
		return false
	}
	if id == constants.Empty_UUID {
		return false
	}
	return true
}
