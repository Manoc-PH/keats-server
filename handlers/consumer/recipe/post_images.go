package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Post_Images(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Images_Req | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Post_Images)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Images_Req | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// Inserting Images
	if insert_recipe_images(db, reqData.Recipe_Image); err != nil {
		log.Println("Post_Images_Req | Error on insert_recipe_images: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	response := schemas.Res_Post_Images{
		Recipe_Image: reqData.Recipe_Image,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func insert_recipe_images(db *sql.DB, recipe_images []schemas.Recipe_Image_Schema) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("insert_recipe_images (Begin) | Error: ", err.Error())
		return err
	}
	// Prepare the SQL statement
	stmt, err := txn.Prepare(
		`INSERT INTO recipe_image (
			recipe_id,
			name_file,
			name_url,
			amount,
			amount_unit,
			amount_unit_desc)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
	)
	if err != nil {
		log.Println("insert_recipe_images (Prepare) | Error: ", err.Error())
		return err
	}
	defer stmt.Close()

	// Insert each row
	for i, img := range recipe_images {
		row := stmt.QueryRow(
			img.Recipe_Id,
			img.Name_File,
			img.Name_URL,
			img.Amount,
			img.Amount_Unit,
			img.Amount_Unit_Desc,
		)
		new_image := schemas.Recipe_Image_Schema{
			Recipe_Id:        img.Recipe_Id,
			Name_File:        img.Name_File,
			Name_URL:         img.Name_URL,
			Amount:           img.Amount,
			Amount_Unit:      img.Amount_Unit,
			Amount_Unit_Desc: img.Amount_Unit_Desc,
		}
		err = row.Scan(&new_image.ID)
		recipe_images[i].ID = new_image.ID
		if err != nil {
			log.Println("insert_recipe_images (Exec) | Error: ", err.Error())
		}
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("insert_recipe_images (commit) | Error: ", err.Error())
		return err
	}
	return nil
}

// Documentation for uploading assets to cloudinary:
// https://cloudinary.com/documentation/upload_images#authenticated_requests