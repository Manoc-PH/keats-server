package utilities

import (
	// "log"
	"os"
	// "github.com/joho/godotenv"
)

func GoDotEnvVariable(key string) string {

	// TODO RESTORE THIS WHEN RUNNING LOCALLY
	// load .env file
	// err := godotenv.Load(".env")

	// if err != nil {
	// 	log.Fatalln("Error loading .env file in helper")
	// }

	return os.Getenv(key)
}
