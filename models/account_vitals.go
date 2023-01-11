package models

import (
	"time"

	"github.com/google/uuid"
)

type Account_Vitals struct {
	ID         uuid.UUID `json:"id"`
	Account_Id uuid.UUID `json:"account_id"`
	// info
	Weight          uint      `json:"weight"`
	Height          uint      `json:"height"`
	Birthday        time.Time `json:"birthday"`
	Sex             string    `json:"sex"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id"`
}
