package schemas

import (
	"time"

	"github.com/google/uuid"
)

// TODO Add birthday and remove age
type Req_Update_Account_Vitals struct {
	Account_ID      uuid.UUID `json:"account_id" validate:"required"`
	Weight          uint      `json:"weight" validate:"required"`
	Height          uint      `json:"height" validate:"required"`
	Birthday        time.Time `json:"birthday" validate:"required"`
	Sex             string    `json:"sex" validate:"required"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id" validate:"required"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id" validate:"required"`
}
