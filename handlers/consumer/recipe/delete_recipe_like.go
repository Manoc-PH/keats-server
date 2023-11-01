package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Delete_Recipe_Like(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Delete_Recipe_Like | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Delete_Recipe_Like)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Delete_Recipe_Like | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// DB transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("Delete_Recipe_Like | Error on db.Begin(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// getting recipe like
	recipe_like := new(models.Recipe_Like)
	err = get_recipe_like(tx, reqData.Recipe_ID, owner_id, recipe_like)
	if err != nil {
		log.Println("Delete_Recipe_Like | Error on get_recipe_like: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// deleting recipe
	err = delete_recipe_like(tx, recipe_like.ID)
	if err != nil {
		log.Println("Delete_Recipe_Like | Error on delete_recipe_like: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// updating recipe
	err = update_recipe_likes(tx, reqData.Recipe_ID, false)
	if err != nil {
		log.Println("Delete_Recipe_Like | Error on update_recipe_likes: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// committing DB transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println("Delete_Recipe_Like | Error on txn.Commit(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}

func get_recipe_like(tx *sql.Tx, recipe_id uuid.UUID, owner_id uuid.UUID, recipe_like *models.Recipe_Like) error {
	row := tx.QueryRow(`SELECT id FROM recipe_like WHERE recipe_id = $1 AND owner_id = $2`, recipe_id, owner_id)
	err := row.Scan(&recipe_like.ID)
	return err
}

func delete_recipe_like(tx *sql.Tx, id uuid.UUID) error {
	_, err := tx.Exec(`DELETE FROM recipe_like WHERE id = $1`, id)
	return err
}
