package handlers

import (
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/food"
	"server/utilities"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/meilisearch/meilisearch-go"
)

func Get_Search_Food(c *fiber.Ctx, db_search *meilisearch.Client) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
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

	formatted_term := strings.Join(strings.Fields(strings.TrimSpace(reqData.Search_Term)), " ")
	res := search_food(db_search, formatted_term)
	return c.Status(fiber.StatusOK).JSON(res)
}

func search_food(db_search *meilisearch.Client, search_term string) []interface{} {
	res, _ := db_search.Index("food").Search(search_term, &meilisearch.SearchRequest{
		AttributesToRetrieve: []string{"id", "name", "name_ph", "name_owner", "barcode", "thumbnail_image_link"},
		MatchingStrategy:     "last",
	})
	return res.Hits
}
