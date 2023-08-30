package handlers

import (
	"database/sql"
	"errors"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/admin/food"
	"server/utilities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Post_Images_Confirm(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Images_Confirm | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	// admin validation
	isAdmin := middlewares.IsAdmin(owner_id, db)
	if isAdmin != true {
		log.Println("Post_Images_Confirm | Error on auth middleware (Not Admin): ")
		return utilities.Send_Error(c, "Only admin users are allowed to access this endpoint", fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Post_Images_Confirm)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Images_Confirm | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// Inserting Images
	if err = confirm_food_images(db, reqData.Food_Images); err != nil {
		log.Println("Post_Images_Confirm | Error on confirm_food_images: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}

func confirm_food_images(db *sql.DB, food_images []models.Food_Image) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("insert_food_images_req (Begin) | Error: ", err.Error())
		return err
	}
	// Prepare the SQL statement
	stmt, err := txn.Prepare(`UPDATE food_image SET name_url = $1 WHERE id = $2`)
	if err != nil {
		log.Println("insert_food_images_req (Prepare) | Error: ", err.Error())
		return err
	}
	defer stmt.Close()

	// Insert each row
	for _, img := range food_images {
		res, err := stmt.Exec(img.Name_URL, img.ID)
		if rows_affected, _ := res.RowsAffected(); rows_affected < 1 {
			log.Println("confirm_food_images (No Rows affected) | Error")
			err = errors.New("Image with id of: " + strconv.Itoa(int(img.ID)) + " not found")
			return err
		}
		if err != nil {
			log.Println("confirm_food_images (Exec) | Error: ", err.Error())
			return err
		}
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("confirm_food_images (commit) | Error: ", err.Error())
		return err
	}
	return nil
}
