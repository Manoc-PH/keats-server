package handlers

import (
	"database/sql"
	"log"
	"net/url"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/setup"
	"server/utilities"
	"strconv"
	"time"

	cld "github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Post_Image(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Images_Req | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Post_Image)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Images_Req | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// Inserting Images
	new_image, err := insert_recipe_image(db, *reqData)
	if err != nil {
		log.Println("Post_Images_Req | Error on insert_recipe_image: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	strTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	values := url.Values{
		"timestamp": []string{strTimestamp},
		"public_id": []string{new_image.Name_File},
	}
	signature, err := cld.SignParameters(
		values,
		setup.Cloudinary_Config.APISecret,
	)

	new_image.Signature = signature
	new_image.Timestamp = strTimestamp
	new_image.Upload_URL = setup.Cloudinary_URL + "/" + setup.Cloudinary_Config.CloudName + "/image/upload"
	new_image.API_key = setup.Cloudinary_Config.APIKey

	return c.Status(fiber.StatusOK).JSON(new_image)
}

func insert_recipe_image(db *sql.DB, recipe_image schemas.Req_Post_Image) (schemas.Res_Post_Image, error) {
	id := uuid.New()
	new_image := schemas.Res_Post_Image{
		ID:             id,
		Recipe_Id:      recipe_image.Recipe_Id,
		Name_File:      "/recipe/images/" + id.String() + ".jpg",
		Name_URL:       setup.Cloudinary_URL + "/" + setup.Cloudinary_Config.CloudName + "/image/upload/recipe/images/" + id.String() + ".jpg",
		Name_URL_Local: recipe_image.Name_URL_Local,
	}
	_, err := db.Exec(`INSERT INTO recipe_image (
		id,
		recipe_id,
		name_file,
		name_url,
		amount,
		amount_unit,
		amount_unit_desc)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		id,
		new_image.Recipe_Id,
		new_image.Name_File,
		new_image.Name_URL,
		new_image.Amount,
		new_image.Amount_Unit,
		new_image.Amount_Unit_Desc)
	if err != nil {
		log.Println("insert_recipe_image | Error: ", err.Error())
		return new_image, err
	}
	return new_image, nil
}

// Documentation for uploading assets to cloudinary:
// https://cloudinary.com/documentation/upload_images#authenticated_requests
