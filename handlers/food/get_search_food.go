package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/food"
	"server/utilities"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Get_Search_Food(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, Owner_Id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Search_Food | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Search_Food)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Search_Food | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	formatted_term := strings.Join(strings.Split(reqData.Search_Term, " "), "<->") + ":*"
	// querying food
	response, err := search_and_scan_food(db, Owner_Id, formatted_term)
	// Server Error
	if err != nil && err != sql.ErrNoRows {
		return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func search_and_scan_food(db *sql.DB, user_id uuid.UUID, search_term string) ([]models.Food, error) {
	rows, err := db.Query(`SELECT
			id,
			name,
			name_ph,
			name_brand,
			food_nutrient_id
		FROM food WHERE search_food @@ to_tsquery($1)
		LIMIT 10;`,
		search_term,
	)
	if err != nil {
		log.Println("Get_Search_Food | error in querying food: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	food_list := make([]models.Food, 0, 10)
	for rows.Next() {
		var new_food = models.Food{}
		if err := rows.
			Scan(
				&new_food.ID,
				&new_food.Name,
				&new_food.Name_Ph,
				&new_food.Name_Brand,
				&new_food.Food_Nutrient_Id,
			); err != nil {
			log.Println("Get_Search_Food | error in scanning food: ", err.Error())
			return nil, err
		}
		food_list = append(food_list, new_food)
	}
	return food_list, err
}
