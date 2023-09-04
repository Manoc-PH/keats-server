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

func Get_Indredients(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Indredients | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	// admin validation
	isAdmin := middlewares.IsAdmin(owner_id, db)
	if isAdmin != true {
		log.Println("Get_Indredients | Error on auth middleware (Not Admin): ")
		return utilities.Send_Error(c, "Only admin users are allowed to access this endpoint", fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Ingredients)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Indredients | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	// querying ingredient
	ingredients := []models.Ingredient{}
	rows, err := query_ingredients(db, reqData.Index, reqData.Size)
	if err != nil && err == sql.ErrNoRows {
		log.Println("Get_Indredients | ingredient does not exist: ", err.Error())
		return utilities.Send_Error(c, "Ingredient does not exist", fiber.StatusBadRequest)
	}
	if err != nil {
		log.Println("error in querying query_ingredients: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		err = scan_ingredients(rows, &ingredients)
		if err != nil && err == sql.ErrNoRows {
			log.Println("Get_Indredients | error in scanning ingredient: ", err.Error())
			return utilities.Send_Error(c, "Ingredient does not exist", fiber.StatusInternalServerError)
		}
		if err != nil {
			log.Println("Get_Indredients | error in scanning ingredient: ", err.Error())
			return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
		}
	}
	return c.Status(fiber.StatusOK).JSON(ingredients)
}

func query_ingredients(db *sql.DB, index uint, size uint) (*sql.Rows, error) {
	rows, err := db.Query(`
		SELECT id, name, name_ph, name_owner, thumbnail_image_link
		FROM ingredient ORDER BY name ASC LIMIT $1 OFFSET $2`,
		size, (size * index),
	)
	return rows, err
}
func scan_ingredients(row *sql.Rows, ingredients *[]models.Ingredient) error {
	var ingredient = models.Ingredient{}
	if err := row.
		Scan(
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.Name_Ph,
			&ingredient.Name_Owner,
			&ingredient.Thumbnail_Image_Link,
		); err != nil {
		return err
	}
	*ingredients = append(*ingredients, ingredient)
	return nil
}
