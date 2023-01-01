package main

import (
	"kryptoverse-api/middlewares"
	"kryptoverse-api/routes"
	"kryptoverse-api/setup"
	"kryptoverse-api/utilities"
	"log"
)

func main() {
	setup.ConnectDB()
	app := setup.SetupApp()
	middlewares.Use_Websocket(app)

	port := utilities.GoDotEnvVariable("PORT")
	if port == "" {
		port = "3000"
	}

	// Http routes
	routes.Auth_Routes(app)

	log.Fatal(app.Listen(":" + port))
}
