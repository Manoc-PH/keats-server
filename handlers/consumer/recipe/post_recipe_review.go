package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Post_Recipe_Review(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Recipe_Review | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Post_Recipe_Review)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Recipe_Review | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// *checking if review exists
	exists := review_exists(owner_id, reqData.Recipe_Id, db)
	if exists == true {
		log.Println("Post_Recipe_Review | Error review exists: ")
		return utilities.Send_Error(c, "You've already made a review, cannot submit another one", fiber.StatusBadRequest)
	}

	//* getting rating and count of recipe
	sum, count, err := get_rating_sum_and_count(db, reqData.Recipe_Id)
	if err != nil {
		log.Println("Post_Recipe_Review | Error on get_rating_sum_and_count: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	new_rating := (float32(sum+int(reqData.Rating)) / float32(count+1))

	//* starting transaction
	txn, err := db.Begin()
	if err != nil {
		log.Println("Post_Recipe_Review | Error on db.Begin(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// updating recipe values
	reqData.Owner_Id = owner_id
	reqData.Date_Created = time.Now()

	// saving recipe review
	err = save_recipe_review(txn, reqData)
	if err != nil {
		log.Println("Post_Recipe_Review | Error on save_recipe_review: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// updating recipe rating
	err = update_recipe_rating(txn, new_rating, reqData.Recipe_Id)

	// committing
	err = txn.Commit()
	if err != nil {
		log.Println("Post_Recipe_Review | Error on txn.Commit(): ", err.Error())
		err = txn.Rollback()
		if err != nil {
			log.Println("Post_Recipe_Review | txn.Rollback(): ", err.Error())
		}
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}
func review_exists(owner_id uuid.UUID, recipe_id uint, db *sql.DB) bool {
	row := db.QueryRow(`
		SELECT id FROM recipe_review 
		WHERE owner_id = $1 AND	recipe_id = $2`, owner_id, recipe_id)
	var id uint
	err := row.Scan(&id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return false
		}
		log.Println("Error in review_exists: ", err.Error())
		return true
	}
	if id != 0 {
		return true
	}
	return false
}
func get_rating_sum_and_count(db *sql.DB, recipe_id uint) (sum int, count int, er error) {
	row := db.QueryRow(`SELECT COUNT(id), COALESCE(SUM(rating), 0) FROM recipe_review WHERE recipe_id = $1`, recipe_id)
	var recipe_count int
	var recipe_sum int
	err := row.Scan(&recipe_count, &recipe_sum)
	if err != nil {
		return 0, 0, err
	}
	return recipe_sum, recipe_count, nil
}
func save_recipe_review(txn *sql.Tx, recipe_review *schemas.Req_Post_Recipe_Review) error {
	row := txn.QueryRow(`INSERT INTO 
		recipe_review(
			description,
			rating,
			owner_id,
			recipe_id,
			date_created)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id`,
		recipe_review.Description,
		recipe_review.Rating,
		recipe_review.Owner_Id,
		recipe_review.Recipe_Id,
		recipe_review.Date_Created,
	)
	err := row.Scan(&recipe_review.ID)
	if err != nil {
		return err
	}
	return nil
}
func update_recipe_rating(txn *sql.Tx, new_rating float32, recipe_id uint) error {
	_, err := txn.Exec(`UPDATE recipe SET rating = $1 WHERE id = $2`, new_rating, recipe_id)
	if err != nil {
		return err
	}
	return nil
}
