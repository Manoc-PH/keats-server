package constants

import "os"

var SecretKey = os.Getenv("SECRET_KEY")
