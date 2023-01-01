package main

import (
	"kryptoverse-api/routes"
	"kryptoverse-api/setup"
	"kryptoverse-api/utilities"
	"log"
)

func main() {
	setup.ConnectDB()
	app := setup.SetupApp()

	port := utilities.GoDotEnvVariable("PORT")
	if port == "" {
		port = "3000"
	}

	// Http routes
	routes.Auth_Routes(app)

	log.Fatal(app.Listen(":" + port))
}
