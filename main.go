package main

import (
	"log"
	"os"
	"server/routes"
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
	routes.Auth_Routes(app)
	routes.Account_Routes(app)
	routes.Tracker_Routes(app)
	routes.Food_Routes(app)

	log.Fatal(app.Listen(":" + port))
}
