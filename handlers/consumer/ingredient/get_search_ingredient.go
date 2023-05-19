package handlers

import (
	"log"
	"server/middlewares"
	schemas "server/schemas/consumer/ingredient"
	"server/utilities"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/meilisearch/meilisearch-go"
)

func Get_Search_Ingredient(c *fiber.Ctx, db_search *meilisearch.Client) error {
	// auth validation
	_, _, err := middlewares.AuthMiddleware(c)
	if err != nil {
		log.Println("Get_Search_Ingredient | Error on auth middleware: ", err.Error())
		return utilities.Send_Error(c, err.Error(), fiber.StatusUnauthorized)
	}
	//* data validation
	reqData := new(schemas.Req_Get_Search_Ingredient)
	if err_data, err := middlewares.Query_Validation(reqData, c); err != nil {
		log.Println("Get_Search_Ingredient | Error on query validation: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err_data)
	}

	formatted_term := strings.Join(strings.Fields(strings.TrimSpace(reqData.Search_Term)), " ")
	log.Println(formatted_term)
	search_ingredient(db_search, formatted_term)
	return nil
	// querying ingredient
	// response, err := search_and_scan_ingredient(db, Owner_Id, formatted_term)
	// Server Error
	// if err != nil && err != sql.ErrNoRows {
	// 	return utilities.Send_Error(c, "An error occured", fiber.StatusInternalServerError)
	// }
	// return c.Status(fiber.StatusOK).JSON(response)
}

func search_ingredient(db_search *meilisearch.Client, search_term string) error {
	res, _ := db_search.Index("ingredients").Search(search_term, &meilisearch.SearchRequest{
		ShowMatchesPosition: true,
	})
	log.Println(res)
	return nil
}
