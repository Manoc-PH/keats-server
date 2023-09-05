package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/admin/ingredient"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Delete_Images(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Delete_Images | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	// admin validation
	isAdmin := middlewares.IsAdmin(owner_id, db)
	if isAdmin != true {
		log.Println("Delete_Images | Error on auth middleware (Not Admin): ")
		return utilities.Send_Error(c, "Only admin users are allowed to access this endpoint", fiber.StatusUnauthorized)
	}
	//* data validation
	req := new(schemas.Req_Delete_Images)
	if err_data, err := middlewares.Body_Validation(req, c); err != nil {
		log.Println("Delete_Images | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// Deleting Images
	err = delete_images_db(db, req.Images)
	if err != nil {
		log.Println("Delete_Images | Error on delete_images_db: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully deleted images",
	})
}

func delete_images_db(db *sql.DB, images []schemas.Ingredient_Image_Schema) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("delete_images_db (Begin) | Error: ", err.Error())
		return err
	}

	// Prepare the SQL statement
	stmt, err := txn.Prepare(`DELETE FROM ingredient_image WHERE id = $1`)
	if err != nil {
		log.Println("delete_images_db (Prepare) | Error: ", err.Error())
		return err
	}
	defer stmt.Close()

	// curl \
	// -d "public_ids[]=image1&public_ids[]=image2" \
	// -X DELETE \
	// https://<API_KEY>:<API_SECRET>@api.cloudinary.com/v1_1/<cloud_name>/resources/image/upload

	// Delete each row
	for _, img := range images {
		// TODO CHECK COUNT OF IMAGES IF NOT DUPLICATE DELETE IN CLOUDINARY
		_, err := stmt.Exec(img.ID)
		if err != nil {
			log.Println("delete_images_db (Exec) | Error: ", err.Error())
			txn.Rollback()
			return err
		}
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("delete_images_db (commit) | Error: ", err.Error())
		return err
	}
	return nil
}
