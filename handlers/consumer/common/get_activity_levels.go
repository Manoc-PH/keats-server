package handlers

import (
	"database/sql"
	"log"
	"server/models"

	"github.com/gofiber/fiber/v2"
)

func Get_Activity_Levels(c *fiber.Ctx, db *sql.DB) error {
	res, err := scan_and_query_activity_levels(db)
	if err != nil {
		log.Println("Error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.JSON(res)
}

func scan_and_query_activity_levels(db *sql.DB) ([]models.Activity_Lvl, error) {
	rows, err := db.Query(`
		SELECT 
			id,
			name,
			coalesce(main_image_link, ''),
			coalesce(background_color, ''),
			coalesce(activity_lvl_desc, '')
		FROM activity_lvl`,
	)
	if err != nil {
		log.Println("Get_Activity_Levels | error in querying Activity_Lvl: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	activity_levels := make([]models.Activity_Lvl, 0, 20)
	for rows.Next() {
		var activity_lvl = models.Activity_Lvl{}
		if err := rows.Scan(
			&activity_lvl.ID,
			&activity_lvl.Name,
			&activity_lvl.Main_Image_Link,
			&activity_lvl.Background_Color,
			&activity_lvl.Activity_Lvl_Desc,
		); err != nil {
			log.Println("Get_Activity_Levels | error in scanning Activity_Lvl: ", err.Error())
			return nil, err
		}
		activity_levels = append(activity_levels, activity_lvl)
	}
	return activity_levels, nil
}
