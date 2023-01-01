package middlewares

import (
	"errors"
	"kryptoverse-api/utilities"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Body_Validation(Req_Struct interface{}, c *fiber.Ctx) error {
	if err := c.BodyParser(Req_Struct); err != nil {
		log.Println("err on line 13: ", err)
		return errors.New("failed validating data")
	}
	err := utilities.ValidateStruct(Req_Struct)
	if err != nil {
		log.Println("err on line 18: ", err)
		return errors.New("failed validating data")
	}
	return nil
}
func Query_Validation(Req_Struct interface{}, c *fiber.Ctx) error {
	if err := c.QueryParser(Req_Struct); err != nil {
		log.Println("err on line 25: ", err)
		return errors.New("failed validating data")
	}
	err := utilities.ValidateStruct(Req_Struct)
	if err != nil {
		log.Println("err on line 30: ", err)
		return errors.New("failed validating data")
	}
	return nil
}
func Struct_Validator(Req_Struct interface{}) error {
	err := utilities.ValidateStruct(Req_Struct)
	if err != nil {
		log.Println("err on line 18: ", err)
		return errors.New("failed validating data")
	}
	return nil
}
