package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Get_Recipe_Review(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Recipe_Review | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := schemas.Req_Get_Recipe_Review{}
	if err_data, err := middlewares.Query_Validation(&reqData, c); err != nil {
		log.Println("Get_Recipe_Review | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// * getting recipe review
	res, err := get_recipe_review(db, reqData, owner_id)
	if err != nil {
		log.Println("Get_Recipe_Review | Error on Get_Recipe_Review: ", err.Error())
		return utilities.Send_Error(c, "An error occured in fetching recipe review", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func get_recipe_review(db *sql.DB, req schemas.Req_Get_Recipe_Review, owner_id uuid.UUID) (schemas.Res_Get_Recipe_Review, error) {
	row := db.QueryRow(`SELECT 
			id,
			description,
			rating,
			owner_id,
			recipe_id,
			date_created
		FROM recipe_review
		WHERE recipe_id = $1 AND owner_id = $2`,
		req.Recipe_Id, owner_id,
	)
	var recipe_rev = schemas.Res_Get_Recipe_Review{}
	if err := row.
		Scan(
			&recipe_rev.ID,
			&recipe_rev.Description,
			&recipe_rev.Rating,
			&recipe_rev.Owner_Id,
			&recipe_rev.Recipe_Id,
			&recipe_rev.Date_Created,
		); err != nil {
		return recipe_rev, err
	}
	return recipe_rev, nil
}
