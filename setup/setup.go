package setup

import (
	"database/sql"
	"kryptoverse-api/utilities"
	"log"
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
	SetupDB(db)
}

func SetupDB(db *sql.DB) error {
	// !TERRIBLE MODEL
	// TOO MUCH DATA DUPLICATION
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id                      SERIAL PRIMARY KEY,
			username                VARCHAR(100) UNIQUE NOT NULL,
			password                VARCHAR(1000) NOT NULL,
			profile_image_link			VARCHAR(1000)
		);
		CREATE TABLE IF NOT EXISTS wallets (
			id                      SERIAL PRIMARY KEY,
			owner_id                INT NOT NULL,
			asset_id     					  INT NOT NULL,
			asset_code              VARCHAR(10) NOT NULL,
			asset_desc              VARCHAR(50) NOT NULL,
			asset_amount            DECIMAL NOT NULL,
			created                 TIMESTAMP,
			FOREIGN KEY(owner_id)   REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS transactions (
			id                      					SERIAL PRIMARY KEY,
			transaction_type_id								INT NOT NULL,
			transaction_type_desc							VARCHAR(100) NOT NULL,
			originator_id           					INT NOT NULL,
			originator_username     					VARCHAR(100) NOT NULL,
			recipient_id            					INT,
			recipient_username	    					VARCHAR(100),
			asset_id													INT NOT NULL,
			asset_code              					VARCHAR(10) NOT NULL,
			asset_desc              					VARCHAR(100) NOT NULL,
			asset_amount            					DECIMAL NOT NULL,
			created														TIMESTAMP,
			FOREIGN KEY(originator_id)      	REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(recipient_id)     		REFERENCES users(id) ON DELETE CASCADE
		); 
		CREATE TABLE IF NOT EXISTS trade_positions (
			id                      					SERIAL PRIMARY KEY,
			position_type_id									INT NOT NULL,
			position_type_desc								VARCHAR(100) NOT NULL,
			position_status_id								INT NOT NULL,
			position_status_desc							VARCHAR(100) NOT NULL,
			originator_id           					INT NOT NULL,
			originator_username     					VARCHAR(100) NOT NULL, 
			asset_id													INT NOT NULL,
			asset_code              					VARCHAR(10) NOT NULL,
			asset_desc              					VARCHAR(50) NOT NULL,
			asset_amount            					DECIMAL NOT NULL,
			buy_asset_price          					DECIMAL NOT NULL,
			buy_asset_amount_usd     					DECIMAL NOT NULL,
			sell_asset_price         					DECIMAL,
			sell_asset_amount_usd    					DECIMAL,
			leverage_amount										INT NOT NULL,
			updated														TIMESTAMP,
			created														TIMESTAMP,
			FOREIGN KEY(originator_id)      	REFERENCES users(id) ON DELETE CASCADE
		); `)
	if err != nil {
		log.Fatalf("an error '%s' was not expected when setting up the db tables", err)
	}
	return nil
}
