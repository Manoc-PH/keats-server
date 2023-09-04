package handlers

import (
	"database/sql"
	"log"
	"net/url"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/admin/food"
	"server/setup"
	"server/utilities"
	"strconv"

	cld "github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/gofiber/fiber/v2"
)

func Post_Images_Req(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Images_Req | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	// admin validation
	isAdmin := middlewares.IsAdmin(owner_id, db)
	if isAdmin != true {
		log.Println("Post_Images_Req | Error on auth middleware (Not Admin): ")
		return utilities.Send_Error(c, "Only admin users are allowed to access this endpoint", fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Post_Images_Req)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Images_Req | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// Inserting Images
	if insert_food_images_req(db, reqData.Food_Images); err != nil {
		log.Println("Post_Images_Req | Error on insert_food_images_req: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	// Generating signature
	strTimestamp := strconv.FormatInt(reqData.Timestamp.Unix(), 10)
	// TODO FIX SIGNATURE GENERATION
	signature, err := cld.SignParameters(url.Values{"timestamp": []string{strTimestamp}}, setup.CloudinaryConfig.APISecret)
	response := schemas.Res_Post_Images_Req{
		Food_Images: reqData.Food_Images,
		Signature:   signature,
		Timestamp:   strTimestamp,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func insert_food_images_req(db *sql.DB, food_images []models.Food_Image) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("insert_food_images_req (Begin) | Error: ", err.Error())
		return err
	}
	// Prepare the SQL statement
	stmt, err := txn.Prepare(
		`INSERT INTO food_image (
			food_id,
			name_file,
			amount,
			amount_unit,
			amount_unit_desc
			)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
	)
	if err != nil {
		log.Println("insert_food_images_req (Prepare) | Error: ", err.Error())
		return err
	}
	defer stmt.Close()

	// Insert each row
	for i, img := range food_images {
		row := stmt.QueryRow(img.Food_Id, img.Name_File, img.Amount, img.Amount_Unit, img.Amount_Unit_Desc)
		new_image := models.Food_Image{
			Food_Id:          img.Food_Id,
			Name_File:        img.Name_File,
			Amount:           img.Amount,
			Amount_Unit:      img.Amount_Unit,
			Amount_Unit_Desc: img.Amount_Unit_Desc,
		}
		err = row.Scan(&new_image.ID)
		food_images[i] = new_image
		if err != nil {
			log.Println("insert_food_images_req (Exec) | Error: ", err.Error())
		}
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("insert_food_images_req (commit) | Error: ", err.Error())
		return err
	}
	return nil
}
