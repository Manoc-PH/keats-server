package main

import (
	"log"
	"os"
	admin_routes "server/routes/admin"
	consumer_routes "server/routes/consumer"
	"server/setup"
)

func main() {
	setup.ConnectDB()
	setup.ConnectAdminDB()
	app := setup.SetupApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Http routes
	// Consumer routes
	consumer_routes.Auth_Routes(app)
	consumer_routes.Account_Routes(app)
	consumer_routes.Tracker_Routes(app)
	consumer_routes.Food_Routes(app)
	consumer_routes.Ingredient_Routes(app)
	consumer_routes.Common_Routes(app)
	// Admin routes
	admin_routes.Ingredient_Routes(app)
	admin_routes.Food_Routes(app)
	admin_routes.Auth_Routes(app)

	log.Fatal(app.Listen(":" + port))
}

// TODO IMPLEMENT DB CON POOLING
// SETUP
// func NewDBPool(connectionString string, maxConnections int) (*sql.DB, error) {
// 	connStr := connectionString

// 	// Configure the connection pool settings
// 	config, err := pq.ParseURL(connStr)
// 	if err != nil {
// 			return nil, err
// 	}
// 	config += fmt.Sprintf(" pool_max_conns=%d", maxConnections)

// 	// Open a new database connection pool
// 	db, err := sql.Open("postgres", config)
// 	if err != nil {
// 			return nil, err
// 	}

// 	// Test the database connection
// 	err = db.Ping()
// 	if err != nil {
// 			db.Close()
// 			return nil, err
// 	}

// 	return db, nil
// }

// MAIN
// const (
// 	dbConnectionString = "your_db_connection_string_here"
// 	maxDBConnections   = 10 // Set your desired maximum number of connections
// )

// func main() {
// 	db, err := NewDBPool(dbConnectionString, maxDBConnections)
// 	if err != nil {
// 			log.Fatalf("Failed to connect to the database: %v", err)
// 	}
// 	defer db.Close()

// 	// Use the 'db' connection pool in your application's handlers
// }

// HANDLER
// func SomeHandler(c *fiber.Ctx) error {
// 	// Get a connection from the pool
// 	conn, err := db.BeginTx(context.Background(), nil)
// 	if err != nil {
// 			return err
// 	}
// 	defer conn.Rollback()

// 	// Execute queries using 'conn'

// 	// Commit the transaction (if successful)
// 	if err := conn.Commit(); err != nil {
// 			return err
// 	}

// 	return nil
// }
