package setup

import (
	"database/sql"
	"log"
	"server/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var SecretKey string
var FiberConfig = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		// Send custom error page
		err = ctx.Status(code).SendFile("./build/notfound.html")
		if err != nil {
			// In case the SendFile fails
			return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		// Return from handler
		return nil
	},
}

func SetupApp() *fiber.App {
	app := fiber.New(FiberConfig)

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	return app
}

func ConnectDB() {
	SecretKey = utilities.GoDotEnvVariable("SECRET_KEY")
	connStr := utilities.GoDotEnvVariable("POSTGRES_URL")
	var err error
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxIdleTime(time.Minute * 2)
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected!")
	DB = db
	// SetupDB(db)
}

func SetupDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS food(
			id 									SERIAL PRIMARY KEY,
			name 								varchar UNIQUE NOT NULL,
			name_ph 						varchar UNIQUE NOT NULL,
			brand_name					varchar,
			barcode							varchar,
			amount 							float4  NOT NULL,
			amount_unit 				varchar(4)  NOT NULL,
			amount_unit_desc 		varchar(40)  NOT NULL,
			serving_size 				float4,
			calories 						float4 NOT NULL,
			protein 						float4 NOT NULL,
			carbs 							float4 NOT NULL,
			fats 								float4 NOT NULL); 
		); `)
	if err != nil {
		log.Fatalf("an error '%s' was not expected when setting up the db tables", err)
	}
	return nil
}
