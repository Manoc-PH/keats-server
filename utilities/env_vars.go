package utilities

import (
	// "log"
	"log"
	"os"

	"github.com/joho/godotenv"
	// "github.com/joho/godotenv"
)

var prod = false

func GoDotEnvVariable(key string) string {

	// load .env file
	if prod == false {
		err := godotenv.Load(".env")

		if err != nil {
			log.Println("Error loading .env file in helper")
			prod = true
		}
	}

	return os.Getenv(key)
}
