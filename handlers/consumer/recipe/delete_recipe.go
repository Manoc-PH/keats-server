package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Delete_Recipe(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Delete_Recipe | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	//* data validation
	reqData := new(schemas.Req_Delete_Recipe)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Delete_Recipe | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// Getting recipe nutrient id
	recipe := new(models.Recipe)
	err = get_recipe_nutrient_id_and_owner(db, reqData.ID, recipe)
	if err != nil {
		log.Println("Delete_Recipe | Error on get_recipe_nutrient_id_and_owner: ", err.Error())
		return utilities.Send_Error(c, "Could not find recipe", fiber.StatusBadRequest)
	}
	if recipe.Owner_Id != owner_id {
		log.Println("Delete_Recipe | Error: attempt of deletion from a user that is not an owner of the recipe")
		return utilities.Send_Error(c, "You cannot delete a recipe that is not yours.", fiber.StatusBadRequest)
	}

	// DB transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("Delete_Recipe | Error on db.Begin(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// deleting recipe
	err = delete_recipe(tx, reqData.ID)
	if err != nil {
		log.Println("Delete_Recipe | Error on delete_recipe: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// deleting recipe nutrient
	err = delete_recipe_nutrient(tx, recipe.Nutrient_Id)
	if err != nil {
		log.Println("Delete_Recipe | Error on delete_recipe_nutrient: ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}

	// committing DB transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println("Delete_Recipe | Error on txn.Commit(): ", err.Error())
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}

func get_recipe_nutrient_id_and_owner(db *sql.DB, recipe_id uint, recipe *models.Recipe) error {
	row := db.QueryRow(`SELECT nutrient_id, owner_id FROM recipe WHERE id = $1`, recipe_id)
	err := row.Scan(&recipe.Nutrient_Id, &recipe.Owner_Id)
	if err != nil {
		return err
	}
	return nil
}
func delete_recipe(tx *sql.Tx, recipe_id uint) error {
	_, err := tx.Exec(`DELETE FROM recipe WHERE id = $1`, recipe_id)
	return err
}
func delete_recipe_nutrient(tx *sql.Tx, nutrient_id uint) error {
	_, err := tx.Exec(`DELETE FROM nutrient WHERE id = $1`, nutrient_id)
	return err
}

// TODO Add function to delete all images from cloudinary and in the db
