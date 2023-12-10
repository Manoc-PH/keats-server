package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/setup"
	"server/utilities"
	"strconv"
	"time"

	cld "github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
)

func Patch_Main_Image(c *fiber.Ctx, db *sql.DB, db_search *meilisearch.Client) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Patch_Main_Image | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Patch_Main_Image)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Patch_Main_Image | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// Inserting Images
	new_image, err := update_recipe_main_image(db, *reqData)
	if err != nil {
		log.Println("Patch_Main_Image | Error on insert_recipe_image: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	strTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	// ! DO NOT USE url.Values IT CRASHES THE SERVER ON CONCCURENT REQS
	image_values := map[string][]string{
		"timestamp": {strTimestamp},
		"public_id": {new_image.Image_Name},
	}
	image_signature, err := cld.SignParameters(image_values, setup.APISecret)
	thumbnail_values := map[string][]string{
		"timestamp": {strTimestamp},
		"public_id": {new_image.Thumbnail_Name},
	}
	thumbnail_signature, err := cld.SignParameters(thumbnail_values, setup.APISecret)

	new_image.Image_Signature = image_signature
	new_image.Thumbnail_Signature = thumbnail_signature
	new_image.Timestamp = strTimestamp
	new_image.Upload_URL = setup.Cloudinary_URL + "/" + setup.CloudName + "/image/upload"
	new_image.API_key = setup.APIKey

	err = update_recipe_image_meili(db_search, new_image)
	if err != nil {
		log.Println("Patch_Main_Image | Error on update_recipe_image_meili: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(new_image)
}

func update_recipe_main_image(db *sql.DB, recipe_image schemas.Req_Patch_Main_Image) (schemas.Res_Patch_Main_Image, error) {
	image_id := uuid.New()
	thumbnail_id := uuid.New()
	new_image := schemas.Res_Patch_Main_Image{
		Recipe_Id:           recipe_image.Recipe_Id,
		Image_Name:          "/recipe/main/images/" + image_id.String(),
		Image_URL:           setup.CLOUDINARY_UPLOADED_URL + "/" + setup.CloudName + "/image/upload/keats/recipe/main/images/" + image_id.String() + ".jpg",
		Image_URL_Local:     recipe_image.Image_URL_Local,
		Thumbnail_Name:      "/recipe/main/images/" + thumbnail_id.String(),
		Thumbnail_URL:       setup.CLOUDINARY_UPLOADED_URL + "/" + setup.CloudName + "/image/upload/keats/recipe/main/images/" + thumbnail_id.String() + ".jpg",
		Thumbnail_URL_Local: recipe_image.Thumbnail_URL_Local,
	}
	_, err := db.Exec(`UPDATE RECIPE SET 
			thumbnail_name = $1,
			thumbnail_url = $2,
			image_name = $3,
			image_url = $4
		WHERE id = $5`,
		new_image.Thumbnail_Name,
		new_image.Thumbnail_URL,
		new_image.Image_Name,
		new_image.Image_URL,
		new_image.Recipe_Id,
	)
	if err != nil {
		log.Println("update_recipe_main_image | Error: ", err.Error())
		return new_image, err
	}
	return new_image, nil
}

// Documentation for uploading assets to cloudinary:
// https://cloudinary.com/documentation/upload_images#authenticated_requests
func update_recipe_image_meili(db_search *meilisearch.Client, recipe schemas.Res_Patch_Main_Image) error {
	new_item := map[string]interface{}{
		"id":        recipe.Recipe_Id,
		"image_url": recipe.Image_URL,
	}
	_, err := db_search.Index("recipes").UpdateDocuments(new_item, "id")
	if err != nil {
		return err
	}
	return nil
}
