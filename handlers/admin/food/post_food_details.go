package handlers

import (
	"database/sql"
	"errors"
	"log"
	"server/middlewares"
	schemas "server/schemas/admin/food"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Post_Food_Details(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Food_Details | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	// admin validation
	isAdmin := middlewares.IsAdmin(owner_id, db)
	if isAdmin != true {
		log.Println("Post_Food_Details | Error on auth middleware (Not Admin): ")
		return utilities.Send_Error(c, "Only admin users are allowed to access this endpoint", fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Post_Food_Details)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Food_Details | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// Inserting data
	if err = insert_food_details(db, reqData); err != nil {
		log.Println("Post_Food_Details | Error on body validation: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}
func insert_food_details(db *sql.DB, food *schemas.Req_Post_Food_Details) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("insert_food_details | Error on starting txn: ", err.Error())
		newErr := errors.New("An error on starting txn: " + err.Error())
		return newErr
	}
	row := txn.QueryRow(`
		INSERT INTO nutrient
			(amount,
			amount_unit,
			amount_unit_desc,
			serving_size,
			calories,
			protein,
			carbs,
			fats,
			trans_fat,
			saturated_fat,
			sugars,
			fiber,
			sodium,
			iron,
			calcium,
			serving_total)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id
	`,
		food.Nutrient.Amount,
		food.Nutrient.Amount_Unit,
		food.Nutrient.Amount_Unit_Desc,
		food.Nutrient.Serving_Size,
		food.Nutrient.Serving_Total,
		food.Nutrient.Calories,
		food.Nutrient.Protein,
		food.Nutrient.Carbs,
		food.Nutrient.Fats,
		food.Nutrient.Trans_Fat,
		food.Nutrient.Saturated_Fat,
		food.Nutrient.Sugars,
		food.Nutrient.Fiber,
		food.Nutrient.Sodium,
		food.Nutrient.Iron,
		food.Nutrient.Calcium)
	err = row.Scan(&food.Nutrient.ID)
	if err != nil {
		log.Println("insert_food_details | Error on saving nutrient: ", err.Error())
		newErr := errors.New("An error occured in saving nutrient: " + err.Error())
		return newErr
	}

	row = txn.QueryRow(`
		INSERT INTO food 
			(name,
			date_created,
			barcode,
			food_desc,
			category_id,
			food_type_id,
			nutrient_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`,
		food.Food.Name,
		time.Now(),
		food.Food.Barcode,
		food.Food.Food_Desc,
		food.Food.Category_Id,
		1,
		food.Nutrient.ID,
	)
	err = row.Scan(&food.Food.ID)
	if err != nil {
		log.Println("insert_food_details | Error on saving food: ", err.Error())
		newErr := errors.New("An error occured in saving food: " + err.Error())
		return newErr
	}

	err = txn.Commit()
	if err != nil {
		log.Println("insert_food_details | Error on commit: ", err.Error())
		txn.Rollback()
		return err
	}
	return nil
}
