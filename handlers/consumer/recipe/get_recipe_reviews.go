package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Get_Recipe_Reviews(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("get_recipe_reviews | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Recipe_Reviews)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("get_recipe_reviews | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	response := new(schemas.Res_Get_Recipe_Reviews)
	// * getting recipe details
	err = get_recipe_reviews(db, *reqData, &response.Reviews)
	if err != nil {
		log.Println("get_recipe_reviews | Error on Get_Recipe_Ingredients: ", err.Error())
		return utilities.Send_Error(c, "An error occured in fetching recipe", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func get_recipe_reviews(db *sql.DB, req schemas.Req_Get_Recipe_Reviews, recipe_revs *[]schemas.Recipe_Reviews_Schema) error {
	rows, err := db.Query(`SELECT 
			recipe_review.id,
			recipe_review.description,
			recipe_review.rating,
			recipe_review.owner_id,
			recipe_review.recipe_id,
			recipe_review.date_created,
			consumer_profile.name_first,
			consumer_profile.name_last
		FROM recipe_review
		JOIN account on recipe_review.owner_id = account.id
		JOIN consumer_profile on account.id = consumer_profile.account_id
		WHERE recipe_review.recipe_id = $1
		LIMIT $2
		OFFSET $3`,
		req.Recipe_Id, req.Size, req.Page*req.Size,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var recipe_rev = schemas.Recipe_Reviews_Schema{}
		var name_first, name_last string
		if err := rows.
			Scan(
				&recipe_rev.ID,
				&recipe_rev.Description,
				&recipe_rev.Rating,
				&recipe_rev.Owner_Id,
				&recipe_rev.Recipe_Id,
				&recipe_rev.Date_Created,
				&name_first,
				&name_last,
			); err != nil {
			return err
		}
		recipe_rev.Name_Owner = name_first + " " + name_last
		*recipe_revs = append(*recipe_revs, recipe_rev)
	}
	return nil
}
