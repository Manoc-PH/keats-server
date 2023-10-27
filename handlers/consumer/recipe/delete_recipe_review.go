package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Delete_Recipe_Review(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Delete_Recipe_Review | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Delete_Recipe_Review)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Delete_Recipe_Review | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// DB transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("Delete_Recipe_Review | Error on db.Begin(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// deleting recipe
	err = delete_recipe_review(tx, reqData.Recipe_ID, owner_id)
	if err != nil {
		log.Println("Delete_Recipe_Review | Error on delete_recipe_like: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// getting rating and count of recipe
	sum, count, err := get_rating_sum_and_count(tx, reqData.Recipe_ID)
	if err != nil {
		log.Println("Post_Recipe_Review | Error on get_rating_sum_and_count: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	new_rating := (float32(sum) / float32(count))

	// updating recipe rating
	err = update_recipe_rating(tx, new_rating, reqData.Recipe_ID, uint(count))
	if err != nil {
		log.Println("Post_Recipe_Review | Error on update_recipe_rating: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// committing DB transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println("Delete_Recipe_Review | Error on txn.Commit(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}

func delete_recipe_review(tx *sql.Tx, recipe_id uint, owner_id uuid.UUID) error {
	_, err := tx.Exec(`DELETE FROM recipe_review WHERE recipe_id = $1 AND owner_id = $2`,
		recipe_id, owner_id)
	return err
}
