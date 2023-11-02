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

	cld "github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/gofiber/fiber/v2"
)

func Post_Thumbnail_Req(c *fiber.Ctx, db *sql.DB) error {
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
	reqData := new(schemas.Req_Post_Thumbnail_Req)
	if err_data, err := middlewares.Body_Validation(reqData, c); err != nil {
		log.Println("Post_Thumbnail_Req | Error on body validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}
	// Generating signature
	strTimestamp := strconv.FormatInt(reqData.Timestamp.Unix(), 10)
	// TODO FIX SIGNATURE GENERATION
	signature, err := cld.SignParameters(url.Values{"timestamp": []string{strTimestamp}}, setup.Cloudinary_Config.APISecret)
	response := schemas.Res_Post_Thumbnail_Req{
		Signature: signature,
		Timestamp: strTimestamp,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
