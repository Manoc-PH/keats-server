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
	consumer_routes.Auth_Routes(app)
	consumer_routes.Account_Routes(app)
	consumer_routes.Tracker_Routes(app)
	consumer_routes.Food_Routes(app)
	consumer_routes.Ingredient_Routes(app)
	consumer_routes.Common_Routes(app)
	admin_routes.Ingredient_Routes(app)

	log.Fatal(app.Listen(":" + port))
}
