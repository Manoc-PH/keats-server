package handlers

import (
	"database/sql"
	"log"
	"server/middlewares"
	schemas "server/schemas/admin/ingredient"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Post_Thumbnail_Confirm(c *fiber.Ctx, db *sql.DB) error {
	// auth validation
	_, owner_id, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Post_Thumbnail_Req | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	// admin validation
	isAdmin := middlewares.IsAdmin(owner_id, db)
	if isAdmin != true {
		log.Println("Post_Thumbnail_Req | Error on auth middleware (Not Admin): ")
		return utilities.Send_Error(c, "Only admin users are allowed to access this endpoint", fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Post_Thumbnail_Confirm)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Thumbnail_Req | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	if err = confirm_thumbnail(db, reqData.Thumbnail_Image_Link, reqData.ID); err != nil {
		log.Println("Post_Thumbnail_Req | Error on confirm_thumbnail: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(reqData)
}

func confirm_thumbnail(db *sql.DB, thumbnail_link string, id uuid.UUID) error {
	_, err := db.Exec(
		`UPDATE ingredient SET thumbnail_image_link = $1 WHERE id = $2`,
		thumbnail_link, id,
	)
	if err != nil {
		log.Println("confirm_thumbnail (Exec) | Error: ", err.Error())
		return err
	}
	return nil
}
