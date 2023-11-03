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
	"github.com/meilisearch/meilisearch-go"
)

// TODO UPDATE MEILISEARCH RECIPE'S RATING AND RATING COUNT
func Post_Recipe_Review(c *fiber.Ctx, db *sql.DB, db_search *meilisearch.Client) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Recipe_Review | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	// data validation
	reqData := new(schemas.Req_Post_Recipe_Review)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Recipe_Review | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// starting transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("Post_Recipe_Review | Error on db.Begin(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// checking if review exists
	exists := review_exists(owner_id, reqData.Recipe_Id, tx)
	if exists == true {
		log.Println("Post_Recipe_Review | Error review exists: ")
		return utilities.Send_Error(c, "You've already made a review, cannot submit another one", fiber.StatusBadRequest)
	}

	// getting rating and count of recipe
	sum, count, err := get_rating_sum_and_count(tx, reqData.Recipe_Id)
	if err != nil {
		log.Println("Post_Recipe_Review | Error on get_rating_sum_and_count: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	new_rating := (float32(sum+int(reqData.Rating)) / float32(count+1))

	// updating recipe values
	reqData.ID = uuid.New()
	reqData.Owner_Id = owner_id
	reqData.Date_Created = time.Now()

	// saving recipe review
	err = save_recipe_review(tx, reqData)
	if err != nil {
		log.Println("Post_Recipe_Review | Error on save_recipe_review: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// updating recipe rating
	err = update_recipe_rating(tx, new_rating, reqData.Recipe_Id, uint(count+1))
	if err != nil {
		log.Println("Post_Recipe_Review | Error on update_recipe_rating: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// updating recipe to meili
	err = update_recipe_rating_meili(db_search, new_rating, reqData.Recipe_Id, uint(count+1))
	if err != nil {
		log.Println("Post_Recipe_Review | Error on update_recipe_rating_meili: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// committing
	err = tx.Commit()
	if err != nil {
		log.Println("Post_Recipe_Review | Error on tx.Commit(): ", err.Error())
		err = tx.Rollback()
		if err != nil {
			log.Println("Post_Recipe_Review | tx.Rollback(): ", err.Error())
		}
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}
func review_exists(owner_id uuid.UUID, recipe_id uuid.UUID, tx *sql.Tx) bool {
	row := tx.QueryRow(`
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
func get_rating_sum_and_count(tx *sql.Tx, recipe_id uuid.UUID) (sum int, count int, er error) {
	row := tx.QueryRow(`SELECT COUNT(id), COALESCE(SUM(rating), 0) FROM recipe_review WHERE recipe_id = $1`, recipe_id)
	var recipe_count int
	var recipe_sum int
	err := row.Scan(&recipe_count, &recipe_sum)
	if err != nil {
		return 0, 0, err
	}
	return recipe_sum, recipe_count, nil
}
func save_recipe_review(tx *sql.Tx, recipe_review *schemas.Req_Post_Recipe_Review) error {
	_, err := tx.Exec(`INSERT INTO 
		recipe_review(
			id,
			description,
			rating,
			owner_id,
			recipe_id,
			date_created)
		VALUES($1, $2, $3, $4, $5, $6)`,
		recipe_review.ID,
		recipe_review.Description,
		recipe_review.Rating,
		recipe_review.Owner_Id,
		recipe_review.Recipe_Id,
		recipe_review.Date_Created,
	)
	if err != nil {
		return err
	}
	return nil
}
func update_recipe_rating(tx *sql.Tx, new_rating float32, recipe_id uuid.UUID, count uint) error {
	_, err := tx.Exec(`UPDATE recipe SET rating = $1, rating_count = $2 WHERE id = $3`, new_rating, count, recipe_id)
	if err != nil {
		return err
	}
	return nil
}
func update_recipe_rating_meili(db_search *meilisearch.Client, new_rating float32, recipe_id uuid.UUID, count uint) error {
	recipe_meili := map[string]interface{}{
		"id":     recipe_id,
		"rating": new_rating,
		"count":  count,
	}
	db_search.Index("recipes").UpdateDocuments(recipe_meili)
	return nil
}
