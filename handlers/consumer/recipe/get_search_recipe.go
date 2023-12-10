package handlers

import (
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/recipe"
	"server/utilities"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/meilisearch/meilisearch-go"
)

func Get_Search_Recipe(c *fiber.Ctx, db_search *meilisearch.Client) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Search_Recipe | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Search_Recipe)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Search_Recipe | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	formatted_term := strings.Join(strings.Fields(strings.TrimSpace(reqData.Search_Term)), " ")
	res := search_recipe(db_search, formatted_term)
	return c.Status(fiber.StatusOK).JSON(res)
}

func search_recipe(db_search *meilisearch.Client, search_term string) []interface{} {
	res, _ := db_search.Index("recipes").Search(search_term, &meilisearch.SearchRequest{
		AttributesToRetrieve: []string{
			"id",
			"name",
			"name_owner",
			"image_url",
			"rating",
			"rating_count",
		},
		MatchingStrategy: "last",
	})
	return res.Hits
}
