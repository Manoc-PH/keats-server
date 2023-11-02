package handlers

import (
	"context"
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/admin/ingredient"
	"server/setup"
	"server/utilities"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/gofiber/fiber/v2"
)

// TODO CREATE ENDPOINT THAT DELETES CLOUDINARY IMAGE IF THIS ENDPOINT FAILS TO DELETE SOME
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
	res, err := delete_images_db(db, req.Images)
	if err != nil {
		log.Println("Delete_Images | Error on delete_images_db: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(*res)
}

func delete_images_db(db *sql.DB, images []schemas.Ingredient_Image_Schema) (*admin.DeleteAssetsResult, error) {
	txn, err := db.Begin()
	if err != nil {
		log.Println("delete_images_db (Begin) | Error: ", err.Error())
		return nil, err
	}

	// Prepare the SQL statement
	stmt, err := txn.Prepare(`DELETE FROM ingredient_image WHERE id = $1`)
	if err != nil {
		log.Println("delete_images_db (Prepare) | Error: ", err.Error())
		return nil, err
	}
	defer stmt.Close()

	// Delete each row
	to_delete_imgs := api.CldAPIArray{}
	for _, img := range images {
		hasDuplicate, err := hasDuplicateImages(db, img.Name_URL)
		if err != nil {
			log.Println("delete_images_db (hasDuplicateImages) | Error: ", err.Error())
			txn.Rollback()
			return nil, err
		}
		if hasDuplicate == false {
			to_delete_imgs = append(to_delete_imgs, img.Name_URL)
		}
		_, err = stmt.Exec(img.ID)
		if err != nil {
			log.Println("delete_images_db (Exec) | Error: ", err.Error())
			txn.Rollback()
			return nil, err
		}
	}
	// Deleting in cloudinary
	res, err := delete_image_cloudinary(to_delete_imgs)
	if err != nil {
		txn.Rollback()
		log.Println("delete_images_db (delete_image_cloudinary) | Error: ", err.Error())
		return nil, err
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("delete_images_db (commit) | Error: ", err.Error())
		return nil, err
	}
	return res, nil
}

func hasDuplicateImages(db *sql.DB, name_url string) (bool, error) {
	count := 0
	row := db.QueryRow(`SELECT count(name_url) FROM ingredient_image WHERE name_url = $1`, name_url)

	err := row.Scan(&count)
	if err != nil {
		return true, err
	}
	if count < 2 {
		return false, nil
	}

	return true, nil
}

// PUBLIC IDS ARE THE URL ENDPOINTS
func delete_image_cloudinary(image_names api.CldAPIArray) (*admin.DeleteAssetsResult, error) {
	apiURL :=
		"https://" +
			setup.Cloudinary_Config.APIKey +
			":" + setup.Cloudinary_Config.APISecret +
			"@api.cloudinary.com/v1_1/" +
			setup.Cloudinary_Config.CloudName +
			"/resources/image/upload"
	cld, err := cloudinary.NewFromURL(apiURL)
	if err != nil {
		return nil, err
	}
	cld.Admin.Config.Cloud.APIKey = setup.Cloudinary_Config.APIKey
	cld.Admin.Config.Cloud.APISecret = setup.Cloudinary_Config.APISecret
	cld.Admin.Config.Cloud.CloudName = setup.Cloudinary_Config.CloudName

	var ctx = context.Background()
	// SAMPLE PUBLIC ID:
	// keats/ingredient/chicken/breast/boneless-raw/logo_filled
	resp, err := cld.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{PublicIDs: image_names})
	if err != nil {
		return resp, err
	}
	return resp, nil

	// apiURL :=
	// 	"https://" +
	// 		setup.CloudinaryConfig.APIKey +
	// 		":" + setup.CloudinaryConfig.APISecret +
	// 		"@api.cloudinary.com/v1_1/" +
	// 		setup.CloudinaryConfig.CloudName +
	// 		"/resources/image/upload"
	// client := &http.Client{}

	// // Create a DELETE request
	// req, err := http.NewRequest("DELETE", apiURL, nil)
	// if err != nil {
	// 	fmt.Println("delete_image_cloudinary | Error creating DELETE request:", err)
	// 	return err
	// }

	// // Send the DELETE request
	// resp, err := client.Do(req)
	// if err != nil {
	// 	fmt.Println("delete_image_cloudinary | Error sending DELETE request:", err)
	// 	return err
	// }
	// defer resp.Body.Close()

	// // Check the response status code
	// if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
	// 	fmt.Println("DELETE request was successful.")
	// } else {
	// 	fmt.Printf("DELETE request failed with status code: %d\n", resp.StatusCode)
	// }
}
