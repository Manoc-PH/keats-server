package handlers

import (
	"database/sql"
	"errors"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Patch_Recipe_Review(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Patch_Recipe_Review | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Patch_Recipe_Review)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Patch_Recipe_Review | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// DB transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("Patch_Recipe_Review | Error on db.Begin(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// updating recipe review
	err = update_recipe_review(tx, reqData)
	if err != nil {
		log.Println("Patch_Recipe_Review | Error on update_recipe_review: ", err.Error())
		return utilities.Send_Error(c, "Recipe Review not found", fiber.StatusBadRequest)
	}

	// getting new rating and count of recipe
	sum, count, err := get_rating_sum_and_count(tx, reqData.Recipe_Id)
	if err != nil {
		log.Println("Patch_Recipe_Review | Error on get_rating_sum_and_count: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	new_rating := (float32(sum+int(reqData.Rating)) / float32(count+1))

	// updating recipe
	err = update_recipe_rating(tx, new_rating, reqData.Recipe_Id, uint(count))
	if err != nil {
		log.Println("Patch_Recipe_Review | Error on update_recipe_rating: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// committing DB transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println("Patch_Recipe_Review | Error on txn.Commit(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}

func update_recipe_review(tx *sql.Tx, review *schemas.Req_Patch_Recipe_Review) error {
	res, err := tx.Exec(`UPDATE recipe_review SET 
			description = $1,
			rating = $2
		WHERE id = $3`, review.Description, review.Rating, review.ID,
	)
	if rows_affected, _ := res.RowsAffected(); rows_affected < 1 {
		return errors.New("Recipe Review not found")
	}
	return err
}
