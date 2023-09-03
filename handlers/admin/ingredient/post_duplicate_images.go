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

func Post_Images_Duplicate(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Duplicate_Images | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	// admin validation
	isAdmin := middlewares.IsAdmin(owner_id, db)
	if isAdmin != true {
		log.Println("Post_Duplicate_Images | Error on auth middleware (Not Admin): ")
		return utilities.Send_Error(c, "Only admin users are allowed to access this endpoint", fiber.StatusUnauthorized)
	}
	//* data validation
	req := new(schemas.Req_Post_Duplicate_Images)
	if err_data, err := middlewares.Body_Validation(req, c); err != nil {
		log.Println("Post_Duplicate_Images | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// Inserting Images
	response := schemas.Res_Post_Duplicate_Images{}
	images, err := duplicate_ingredient_images(db, req.Ingredient_Mapping_Id, req.Copied_Ingredient_Mapping_Id)
	if err != nil {
		log.Println("Post_Duplicate_Images | Error on insert_ingredient_images: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	response.Ingredient_Images = images
	return c.Status(fiber.StatusOK).JSON(response)
}

func duplicate_ingredient_images(db *sql.DB, mapping_id uint, copied_mapping_id uint) ([]models.Ingredient_Image, error) {
	// Querying existing images
	rows, err := db.Query(`SELECT 
			ingredient_mapping_id,
			name_file,
			amount,
			amount_unit,
			amount_unit_desc,
			name_url 
		FROM ingredient_image
		WHERE ingredient_mapping_id = $1
	`, copied_mapping_id)
	if err != nil {
		log.Println("insert_ingredient_images (Query) | Error: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	new_images := []models.Ingredient_Image{}
	for rows.Next() {
		img := models.Ingredient_Image{}
		err := rows.Scan(&img.Ingredient_Mapping_Id, &img.Name_File, &img.Amount, &img.Amount_Unit, &img.Amount_Unit_Desc, &img.Name_URL)
		img.Ingredient_Mapping_Id = mapping_id
		if err != nil {
			log.Println("insert_ingredient_images (Query Scan) | Error: ", err.Error())
			return nil, err
		}
		new_images = append(new_images, img)
	}

	txn, err := db.Begin()
	if err != nil {
		log.Println("insert_ingredient_images (Begin) | Error: ", err.Error())
		return nil, err
	}

	// Prepare the SQL statement
	stmt, err := txn.Prepare(
		`INSERT INTO ingredient_image (
				ingredient_mapping_id,
				name_file,
				amount,
				amount_unit,
				amount_unit_desc,
				name_url
			)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
	)
	if err != nil {
		log.Println("insert_ingredient_images (Prepare) | Error: ", err.Error())
		return nil, err
	}
	defer stmt.Close()

	// Insert each row
	for i, img := range new_images {
		row := stmt.QueryRow(img.Ingredient_Mapping_Id, img.Name_File, img.Amount, img.Amount_Unit, img.Amount_Unit_Desc, "")
		new_image := models.Ingredient_Image{
			Ingredient_Mapping_Id: img.Ingredient_Mapping_Id,
			Name_File:             img.Name_File,
			Amount:                img.Amount,
			Amount_Unit:           img.Amount_Unit,
			Amount_Unit_Desc:      img.Amount_Unit_Desc,
		}
		err = row.Scan(&new_image.ID)
		new_images[i] = new_image
		if err != nil {
			log.Println("insert_ingredient_images (Exec) | Error: ", err.Error())
		}
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("insert_ingredient_images (commit) | Error: ", err.Error())
		return nil, err
	}
	return new_images, nil
}
