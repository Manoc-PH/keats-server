package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Req_Sign_Up struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=10,max=32"` //,missingRequiredCharacters // add this for password validation

	Name_First      string    `json:"name_first" validate:"required"`
	Name_Last       string    `json:"name_last" validate:"required"`
	Weight          int       `json:"weight" validate:"required,min=1,max=200"`
	Height          int       `json:"height" validate:"required,min=1,max=250"`
	Birthday        time.Time `json:"birthday" validate:"required"`
	Sex             string    `json:"sex" validate:"required,min=1,max=1"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id" validate:"required,min=1,max=32"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id" validate:"required,min=1,max=32"`
	// time.Time SHOULD BE IN ISO STRING
}
type Res_Sign_Up struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	Weight          int       `json:"weight"`
	Height          int       `json:"height"`
	Birthday        time.Time `json:"birthday"`
	Sex             string    `json:"sex"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id"`
	Token           string    `json:"token"`
	Account_Type_Id uuid.UUID `json:"account_type_id"`
}
