package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/admin/ingredient"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Post_Ingredient_Images_Confirm(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Ingredient_Images_Confirm | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	// admin validation
	isAdmin := middlewares.IsAdmin(owner_id, db)
	if isAdmin != true {
		log.Println("Post_Ingredient_Images_Confirm | Error on auth middleware (Not Admin): ")
		return utilities.Send_Error(c, "Only admin users are allowed to access this endpoint", fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Post_Ingredient_Images_Confirm)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Ingredient_Images_Confirm | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// Inserting Images
	if err = confirm_ingredient_images(db, reqData.Ingredient_Images); err != nil {
		log.Println("Post_Ingredient_Images_Confirm | Error on confirm_ingredient_images: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}

func confirm_ingredient_images(db *sql.DB, ingredient_images []models.Ingredient_Image) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("insert_ingredient_images_req (Begin) | Error: ", err.Error())
		return err
	}
	// Prepare the SQL statement
	stmt, err := txn.Prepare(`UPDATE ingredient_image SET name_url = $1 WHERE id = $2`)
	if err != nil {
		log.Println("insert_ingredient_images_req (Prepare) | Error: ", err.Error())
		return err
	}
	defer stmt.Close()

	// Insert each row
	for _, img := range ingredient_images {
		_, err := stmt.Exec(img.Name_File, img.ID)
		if err != nil {
			log.Println("confirm_ingredient_images (Exec) | Error: ", err.Error())
		}
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("confirm_ingredient_images (commit) | Error: ", err.Error())
		return err
	}
	return nil
}