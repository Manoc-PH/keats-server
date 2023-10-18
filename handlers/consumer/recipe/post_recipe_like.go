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

func Post_Recipe_Like(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Recipe_Like | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Post_Recipe_Like)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Recipe_Like | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// *checking if like exists
	exists := like_exists(owner_id, reqData.Recipe_Id, db)
	if exists == true {
		log.Println("Post_Recipe_Like | Error like exists: ")
		return utilities.Send_Error(c, "You've already liked this recipe.", fiber.StatusBadRequest)
	}

	//* starting transaction
	txn, err := db.Begin()
	if err != nil {
		log.Println("Post_Recipe_Like | Error on db.Begin(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// saving recipe like
	err = save_recipe_like(txn, reqData.Recipe_Id, owner_id)
	if err != nil {
		log.Println("Post_Recipe_Like | Error on save_recipe_like: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	err = update_recipe_likes(txn, reqData.Recipe_Id)
	if err != nil {
		log.Println("Post_Recipe_Like | Error on update_recipe_likes: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// committing
	err = txn.Commit()
	if err != nil {
		log.Println("Post_Recipe_Like | Error on txn.Commit(): ", err.Error())
		err = txn.Rollback()
		if err != nil {
			log.Println("Post_Recipe_Like | txn.Rollback(): ", err.Error())
		}
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}
func like_exists(owner_id uuid.UUID, recipe_id uint, db *sql.DB) bool {
	row := db.QueryRow(`
		SELECT id FROM recipe_like 
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
func save_recipe_like(txn *sql.Tx, recipe_id uint, owner_id uuid.UUID) error {
	_, err := txn.Exec(`INSERT INTO recipe_like
			(owner_id, date_created, recipe_id)
			VALUES($1, $2, $3) `,
		owner_id,
		time.Now(),
		recipe_id,
	)
	if err != nil {
		return err
	}
	return nil
}
func update_recipe_likes(txn *sql.Tx, recipe_id uint) error {
	_, err := txn.Exec(`UPDATE recipe SET likes = likes + 1 WHERE id = $1`,
		recipe_id,
	)
	if err != nil {
		return err
	}
	return nil
}
