package handlers

import (
	"database/sql"
	"log"
	"server/models"

	"github.com/gofiber/fiber/v2"
)

func Get_Diet_Plans(c *fiber.Ctx, db *sql.DB) error {
	res, err := scan_and_query_diet_plans(db)
	if err != nil {
		log.Println("Error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.JSON(res)
}

func scan_and_query_diet_plans(db *sql.DB) ([]models.Diet_Plan, error) {
	rows, err := db.Query(`
	SELECT 
		id,
		name,
		coalesce(main_image_link, ''),
		coalesce(background_color, ''),
		coalesce(diet_plan_desc, '')
	FROM diet_plan
	ORDER BY calorie_percentage DESC`,
	)
	if err != nil {
		log.Println("Get_Diet_Plans | error in querying Diet_Plan: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	diet_plans := make([]models.Diet_Plan, 0, 20)
	for rows.Next() {
		var diet_plan = models.Diet_Plan{}
		if err := rows.Scan(
			&diet_plan.ID,
			&diet_plan.Name,
			&diet_plan.Main_Image_Link,
			&diet_plan.Background_Color,
			&diet_plan.Diet_Plan_Desc,
		); err != nil {
			log.Println("Get_Diet_Plans | error in scanning Diet_Plan: ", err.Error())
			return nil, err
		}
		diet_plans = append(diet_plans, diet_plan)
	}
	return diet_plans, nil
}
