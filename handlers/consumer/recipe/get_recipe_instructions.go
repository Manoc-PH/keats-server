package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	"server/models"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

func Get_Recipe_Instructions(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("get_recipe_instructions | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Recipe_Instructions)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("get_recipe_instructions | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	response := new(schemas.Res_Get_Recipe_Instructions)
	// * getting recipe details
	err = get_recipe_instructions(db, reqData.Recipe_Id, &response.Instructions)
	if err != nil {
		log.Println("get_recipe_instructions | Error on Get_Recipe_Ingredients: ", err.Error())
		return utilities.Send_Error(c, "An error occured in fetching recipe", fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func get_recipe_instructions(db *sql.DB, recipe_id uint, recipe_ings *[]models.Recipe_Instruction) error {
	rows, err := db.Query(`SELECT 
			id,
			recipe_id,
			instruction_description,
			step_num
		FROM recipe_instruction 
		WHERE recipe_id = $1`,
		recipe_id,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var recipe_inst = models.Recipe_Instruction{}
		if err := rows.
			Scan(
				&recipe_inst.ID,
				&recipe_inst.Recipe_Id,
				&recipe_inst.Instruction_Description,
				&recipe_inst.Step_Num,
			); err != nil {
			return err
		}
		*recipe_ings = append(*recipe_ings, recipe_inst)
	}
	return nil
}
