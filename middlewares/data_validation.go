package middlewares

import (
	"errors"
	"log"
	"server/utilities"

	"github.com/gofiber/fiber/v2"
)

type Error_Data struct {
	Message string                     `json:"message"`
	Data    []*utilities.ErrorResponse `json:"data"`
}

func Body_Validation(Req_Struct interface{}, c *fiber.Ctx) (Error_Data, error) {
	if err := c.BodyParser(Req_Struct); err != nil {
		err_data := Error_Data{Message: "Error in parsing body", Data: nil}
		return err_data, err
	}
	err := utilities.ValidateStruct(Req_Struct)
	if err != nil {
		err_data := Error_Data{Message: "missing or invalid data sent", Data: err}
		return err_data, errors.New("missing or invalid data sent")
	}
	return Error_Data{}, nil
}
func Query_Validation(Req_Struct interface{}, c *fiber.Ctx) (Error_Data, error) {
	if err := c.QueryParser(Req_Struct); err != nil {
		err_data := Error_Data{Message: "Error in parsing query", Data: nil}
		return err_data, err
	}
	err := utilities.ValidateStruct(Req_Struct)
	if err != nil {
		err_data := Error_Data{Message: "Error in parsing query", Data: err}
		return err_data, errors.New("missing or invalid data sent")
	}
	return Error_Data{}, nil
}
func Struct_Validator(Req_Struct interface{}) error {
	err := utilities.ValidateStruct(Req_Struct)
	if err != nil {
		log.Println("err on line 18: ", err)
		return errors.New("failed validating data")
	}
	return nil
}
