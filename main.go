package main

import (
	"log"
	"os"
	consumer_routes "server/routes/consumer"
	"server/setup"
)

func main() {
	setup.ConnectDB()
	app := setup.SetupApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Http routes
	consumer_routes.Auth_Routes(app)
	consumer_routes.Account_Routes(app)
	consumer_routes.Tracker_Routes(app)
	consumer_routes.Food_Routes(app)

	log.Fatal(app.Listen(":" + port))
}
