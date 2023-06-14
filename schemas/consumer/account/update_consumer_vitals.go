package schemas

import (
	"server/models"
	"time"

	"github.com/google/uuid"
)

type Req_Update_Consumer_Vitals struct {
	Account_ID      uuid.UUID `json:"account_id" validate:"required"`
	Weight          uint      `json:"weight" validate:"required"`
	Height          uint      `json:"height" validate:"required"`
	Birthday        time.Time `json:"birthday" validate:"required"`
	Sex             string    `json:"sex" validate:"required"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id" validate:"required"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id" validate:"required"`
}
type Res_Update_Consumer_Vitals struct {
	ReqData         Req_Update_Consumer_Vitals `json:"consumer_vitals"`
	Daily_Nutrients models.Daily_Nutrients     `json:"daily_nutrients"`
}
