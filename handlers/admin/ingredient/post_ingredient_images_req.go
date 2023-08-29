package handlers

import (
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/admin/ingredient"
	"server/setup"
	"server/utilities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Post_Ingredient_Images_Req(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Ingredient_Images_Req | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	// admin validation
	isAdmin := middlewares.IsAdmin(owner_id, db)
	if isAdmin != true {
		log.Println("Post_Ingredient_Images_Req | Error on auth middleware (Not Admin): ")
		return utilities.Send_Error(c, "Only admin users are allowed to access this endpoint", fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Post_Ingredient_Images)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Ingredient_Images_Req | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// Inserting Images
	insert_ingredient_images_req(db, reqData.Ingredient_Images)
	// Generating signature
	strTimestamp := strconv.FormatInt(reqData.Timestamp.Unix(), 10)
	// TODO FIX SIGNATURE GENERATION
	signature := generateCloudinarySignature(setup.CloudinaryConfig.APISecret, strTimestamp, setup.CloudinaryConfig.CloudName)
	response := schemas.Res_Post_Ingredient_Images{
		Ingredient_Images: reqData.Ingredient_Images,
		Signature:         signature,
		Timestamp:         reqData.Timestamp.Unix(),
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func insert_ingredient_images_req(db *sql.DB, ingredient_images []models.Ingredient_Image) error {
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	defer stmt.Close()

	// Insert each row
	for i, img := range ingredient_images {
		row := stmt.QueryRow(img.Ingredient_Mapping_Id, img.Name_File, img.Amount, img.Amount_Unit, img.Amount_Unit_Desc, "")
		new_image := models.Ingredient_Image{
			Ingredient_Mapping_Id: img.Ingredient_Mapping_Id,
			Name_File:             img.Name_File,
			Amount:                img.Amount,
			Amount_Unit:           img.Amount_Unit,
			Amount_Unit_Desc:      img.Amount_Unit_Desc,
		}
		err = row.Scan(&new_image.ID)
		ingredient_images[i] = new_image
		if err != nil {
			log.Println("insert_ingredient_images_req (commit) | Error: ", err.Error())
		}
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("insert_ingredient_images_req (commit) | Error: ", err.Error())
		return err
	}
	return nil
}
func generateCloudinarySignature(apiSecret string, timestamp string, publicID string) string {
	toSign := "&timestamp=" + timestamp + apiSecret
	h := hmac.New(sha1.New, []byte(apiSecret))
	h.Write([]byte(toSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
