package setup

import (
	"database/sql"
	"fmt"
	"log"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"github.com/meilisearch/meilisearch-go"
)

var DB *sql.DB
var Admin_DB *sql.DB
var DB_Search *meilisearch.Client
var SecretKey string
var FiberConfig = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		// Send custom error page
		if code == fiber.StatusMethodNotAllowed {
			return ctx.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
				"message": "Method not allowed",
			})
		}
		if code == fiber.StatusNotFound {
			return ctx.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
				"message": "Route does not exist",
			})
		}
		return ctx.Status(code).JSON(fiber.Map{
			"message": err.Error(),
		})
	},
}

// Cloudinary API credentials
type CloudinaryConfigType struct {
	CloudName string
	APIKey    string
	APISecret string
}

var CloudinaryConfig CloudinaryConfigType

func SetupApp() *fiber.App {
	app := fiber.New(FiberConfig)

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	CloudName := utilities.GoDotEnvVariable("CLOUDINARY_NAME")
	APIKey := utilities.GoDotEnvVariable("CLOUDINARY_API_KEY")
	APISecret := utilities.GoDotEnvVariable("CLOUDINARY_API_SECRET")
	CloudinaryConfig = CloudinaryConfigType{
		CloudName: CloudName,
		APIKey:    APIKey,
		APISecret: APISecret,
	}

	return app
}

func ConnectDB() {
	SecretKey = utilities.GoDotEnvVariable("SECRET_KEY")
	// connStr := utilities.GoDotEnvVariable("POSTGRES_URL")
	dbuser := utilities.GoDotEnvVariable("DB_USER")
	dbpass := utilities.GoDotEnvVariable("DB_PASSWORD")
	dbname := utilities.GoDotEnvVariable("DB_DB")
	dbhost := utilities.GoDotEnvVariable("DB_HOST")
	var err error
	psqlInfo := fmt.Sprintf(`host=%s port=%d user=%v password=%v dbname=%v sslmode=disable`,
		dbhost, 5432, dbuser, dbpass, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxIdleTime(time.Minute * 2)
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected to Postgres!")
	DB = db

	db_search_api_key := utilities.GoDotEnvVariable("MEILISEARCH_ADMIN_KEY")
	meilisearch_host := utilities.GoDotEnvVariable("MEILISEARCH_HOST")
	// !WHEN RUNNING ON DOCKER CHANGE THE HOST TO THE CONTAINER NAME
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   meilisearch_host,
		APIKey: db_search_api_key,
	})
	if client != nil {
		log.Println("Connected to Meilisearch!")
	}
	DB_Search = client
	setupMeiliIngredients(db, client)
}
func ConnectAdminDB() {
	dbuser := utilities.GoDotEnvVariable("ADMIN_DB_USER")
	dbpass := utilities.GoDotEnvVariable("ADMIN_DB_PASSWORD")
	dbname := utilities.GoDotEnvVariable("ADMIN_DB_DB")
	dbhost := utilities.GoDotEnvVariable("ADMIN_DB_HOST")
	var err error
	psqlInfo := fmt.Sprintf(`host=%s port=%d user=%v password=%v dbname=%v sslmode=disable`,
		dbhost, 5432, dbuser, dbpass, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("Error trying to connect to admin db account: ", err)
	}
	db.SetConnMaxIdleTime(time.Minute * 2)
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected to Admin Postgres!")
	Admin_DB = db
}

// func SetupDB(db *sql.DB) error {
// 	_, err := db.Exec(`
// 		CREATE TABLE IF NOT EXISTS food(
// 			id 									SERIAL PRIMARY KEY,
// 			name 								varchar UNIQUE NOT NULL,
// 			name_ph 						varchar UNIQUE NOT NULL,
// 			brand_name					varchar,
// 			barcode							varchar,
// 			amount 							float4  NOT NULL,
// 			amount_unit 				varchar(4)  NOT NULL,
// 			amount_unit_desc 		varchar(40)  NOT NULL,
// 			serving_size 				float4,
// 			calories 						float4 NOT NULL,
// 			protein 						float4 NOT NULL,
// 			carbs 							float4 NOT NULL,
// 			fats 								float4 NOT NULL);
// 		); `)
// 	if err != nil {
// 		log.Fatalf("an error '%s' was not expected when setting up the db tables", err)
// 	}
// 	return nil
// }

func setupMeiliIngredients(db *sql.DB, db_search *meilisearch.Client) error {
	numOfRows := 0
	row := db.QueryRow(`SELECT COUNT(name) FROM ingredient`)
	err := row.Scan(&numOfRows)
	if err != nil {
		return nil
	}
	meili_stats, err := db_search.GetStats()
	if err != nil {
		log.Println(err)
		log.Panicln("Could not get stats of meili db")
	}
	if meili_stats.Indexes["ingredients"].NumberOfDocuments != int64(numOfRows) {
		db_search.Index("ingredients").DeleteAllDocuments()
		_, err = db_search.CreateIndex(&meilisearch.IndexConfig{
			Uid:        "ingredients",
			PrimaryKey: "id",
		})
		filterableAttributes := []string{
			"name",
		}
		db_search.Index("ingredients").UpdateFilterableAttributes(&filterableAttributes)
		if err != nil {
			log.Panicln("Could not create index for ingredient in meili db")
		}
		insert_ingredients(db, db_search)
	}
	return nil
}
func insert_ingredients(db *sql.DB, db_search *meilisearch.Client) {
	type ingredient_mapping struct {
		Ingredient_Mapping_Id      uint   `json:"ingredient_mapping_id"`
		Thumbnail_Image_Link       string `json:"thumbnail_image_link"`
		Ingredient_Id              uint   `json:"ingredient_id"`
		Ingredient_Name            string `json:"ingredient_name"`
		Ingredient_Name_Ph         string `json:"ingredient_name_ph"`
		Ingredient_Name_Owner      string `json:"ingredient_name_owner"`
		Ingredient_Variant_Id      uint   `json:"ingredient_variant_id"`
		Ingredient_Variant_Name    string `json:"ingredient_variant_name"`
		Ingredient_Subvariant_Id   uint   `json:"ingredient_subvariant_id"`
		Ingredient_Subvariant_Name string `json:"ingredient_subvariant_name"`
	}
	type ingredient_details struct {
		Ingredient_Mapping_Id      uint   `json:"ingredient_mapping_id"`
		Ingredient_Variant_Name    string `json:"ingredient_variant_name"`
		Ingredient_Subvariant_Name string `json:"ingredient_subvariant_name"`
	}
	// *This structure works for meilisearch
	// Using showMatchesPosition parameter when searching we can find the match
	// inside the array of ingredients. More information here:
	// https://www.meilisearch.com/docs/reference/api/search#show-matches-position
	// TODO ADD NAME_PH
	type edible struct {
		Id                   uint                 `json:"ingredient_id"`
		Name                 string               `json:"name"`
		Name_Ph              string               `json:"name_ph"`
		Name_Owner           string               `json:"name_owner"`
		Thumbnail_Image_Link string               `json:"thumbnail_image_link"`
		Ingredient_Details   []ingredient_details `json:"ingredient_details"`
	}
	docs := map[string]edible{}
	rows, err := db.Query(`
	SELECT 
		ingredient_mapping.id,
		ingredient.id,
		ingredient.name,
		ingredient.name_ph,
		ingredient.name_owner,
		ingredient.thumbnail_image_link,
		coalesce(ingredient_variant.id, 0),
		coalesce(ingredient_variant.name, ''),
		coalesce(ingredient_subvariant.id, 0),
		coalesce(ingredient_subvariant.name, '')
	FROM ingredient_mapping
	JOIN ingredient on ingredient_mapping.ingredient_id = ingredient.id
	JOIN ingredient_variant on ingredient_mapping.ingredient_variant_id = ingredient_variant.id
	JOIN ingredient_subvariant on ingredient_mapping.ingredient_subvariant_id = ingredient_subvariant.id`)
	if err != nil {
		log.Println("Error querying ingredient: ", err.Error())
	}
	for rows.Next() {
		var new_ing = ingredient_mapping{}
		if err := rows.
			Scan(
				&new_ing.Ingredient_Mapping_Id,
				&new_ing.Ingredient_Id,
				&new_ing.Ingredient_Name,
				&new_ing.Ingredient_Name_Ph,
				&new_ing.Ingredient_Name_Owner,
				&new_ing.Thumbnail_Image_Link,
				&new_ing.Ingredient_Variant_Id,
				&new_ing.Ingredient_Variant_Name,
				&new_ing.Ingredient_Subvariant_Id,
				&new_ing.Ingredient_Subvariant_Name,
			); err != nil {
			log.Println("Error scanning ingredient: ", err.Error())
		}
		var new_ing_details = ingredient_details{
			Ingredient_Mapping_Id:      new_ing.Ingredient_Mapping_Id,
			Ingredient_Variant_Name:    new_ing.Ingredient_Variant_Name,
			Ingredient_Subvariant_Name: new_ing.Ingredient_Subvariant_Name,
		}
		if entry, ok := docs[new_ing.Ingredient_Name]; ok {
			entry.Ingredient_Details = append(entry.Ingredient_Details, new_ing_details)
			docs[new_ing.Ingredient_Name] = entry
		} else {
			new_edible := edible{
				Id:                   new_ing.Ingredient_Id,
				Name:                 new_ing.Ingredient_Name,
				Name_Ph:              new_ing.Ingredient_Name_Ph,
				Name_Owner:           new_ing.Ingredient_Name_Owner,
				Thumbnail_Image_Link: new_ing.Thumbnail_Image_Link,
			}
			new_edible.Ingredient_Details = append(new_edible.Ingredient_Details, new_ing_details)
			docs[new_ing.Ingredient_Name] = new_edible
		}
	}
	formatted_doc := []map[string]interface{}{}
	for _, item := range docs {
		new_item := []map[string]interface{}{{
			"id":                   item.Id,
			"name":                 item.Name,
			"name_ph":              item.Name_Ph,
			"name_owner":           item.Name_Owner,
			"thumbnail_image_link": item.Thumbnail_Image_Link,
			"ingredient_details":   item.Ingredient_Details,
		}}
		formatted_doc = append(formatted_doc, new_item[0])
	}
	_, err = db_search.Index("ingredients").AddDocuments(formatted_doc)
	if err == nil {
		log.Println("Successfully added ingredients to Meilisearch")
	}
}
