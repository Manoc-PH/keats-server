package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Req_Sign_Up struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=10,max=32"` //,missingRequiredCharacters // add this for password validation

	Weight          int       `json:"weight" validate:"required,min=1,max=200"`
	Height          int       `json:"height" validate:"required,min=1,max=250"`
	Birthday        time.Time `json:"birthday" validate:"required"`
	Sex             string    `json:"sex" validate:"required,min=1,max=1"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id" validate:"required,min=1,max=32"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id" validate:"required,min=1,max=32"`
	Measure_Unit_Id uuid.UUID `json:"measure_unit_id" validate:"required,min=1,max=32"`
	// time.Time SHOULD BE IN ISO STRING
}
