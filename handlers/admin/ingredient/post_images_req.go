package handlers

import (
	"database/sql"
	"log"
	"net/url"
	"server/middlewares"
	schemas "server/schemas/admin/ingredient"
	"server/setup"
	"server/utilities"
	"strconv"
	"strings"

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
	if err = insert_ingredient_images_req(db, reqData.Ingredient_Images); err != nil {
		log.Println("Post_Images_Req | Error on insert_ingredient_images_req: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	// Generating signature
	strTimestamp := strconv.FormatInt(reqData.Timestamp.Unix(), 10)
	// TODO FIX SIGNATURE GENERATION
	signature, err := cld.SignParameters(url.Values{"timestamp": []string{strTimestamp}}, setup.CloudinaryConfig.APISecret)
	response := schemas.Res_Post_Images_Req{
		Ingredient_Images: reqData.Ingredient_Images,
		Signature:         signature,
		Timestamp:         strTimestamp,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func insert_ingredient_images_req(db *sql.DB, ingredient_images []schemas.Ingredient_Image_Req) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("insert_ingredient_images_req (Begin) | Error: ", err.Error())
		return err
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
		log.Println("insert_ingredient_images_req (Prepare) | Error: ", err.Error())
		return err
	}
	defer stmt.Close()

	// Insert each row
	for i, img := range ingredient_images {
		name_file, err := generate_ingredient_img_name(db, img.Ingredient_Mapping_Id, img.Amount)
		if err != nil {
			log.Println("insert_ingredient_images_req (generate_ingredient_img_name) | Error: ", err.Error())
			txn.Rollback()
			return err
		}
		row := stmt.QueryRow(img.Ingredient_Mapping_Id, name_file, img.Amount, img.Amount_Unit, img.Amount_Unit_Desc, "")
		new_image := schemas.Ingredient_Image_Req{
			Ingredient_Mapping_Id: img.Ingredient_Mapping_Id,
			Name_File:             name_file,
			Amount:                img.Amount,
			Amount_Unit:           img.Amount_Unit,
			Amount_Unit_Desc:      img.Amount_Unit_Desc,
		}
		err = row.Scan(&new_image.ID)
		if err != nil {
			log.Println("insert_ingredient_images_req (Exec) | Error: ", err.Error())
			txn.Rollback()
			return err
		}
		ingredient_images[i] = new_image
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		log.Println("insert_ingredient_images_req (commit) | Error: ", err.Error())
		return err
	}
	return nil
}
func generate_ingredient_img_name(db *sql.DB, ingredient_mapping_id uint, amount float32) (string, error) {
	name := ""
	variant_name := ""
	subvariant_name := ""
	row := db.QueryRow(`
		SELECT ingredient.name, ingredient_variant.name, ingredient_subvariant.name
		FROM ingredient_mapping
		JOIN ingredient ON ingredient_mapping.ingredient_id = ingredient.id
		JOIN ingredient_variant ON ingredient_mapping.ingredient_variant_id = ingredient_variant.id
		JOIN ingredient_subvariant ON ingredient_mapping.ingredient_subvariant_id = ingredient_subvariant.id
		WHERE ingredient_mapping.id = $1
	`, ingredient_mapping_id)

	err := row.Scan(&name, &variant_name, &subvariant_name)
	if err != nil {
		return "", err
	}
	name = strings.Join(strings.Split(name, " "), "_")
	variant_name = strings.Join(strings.Split(variant_name, " "), "_")
	subvariant_name = strings.Join(strings.Split(subvariant_name, " "), "_")
	amount_str := strconv.FormatFloat(float64(amount), 'f', -1, 32)
	finalData := "keats/ingredient/" + name + "/" + variant_name + "/" + subvariant_name + "/" + amount_str
	return finalData, nil
}
