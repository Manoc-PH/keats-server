package schemas

import (
	"time"

	"github.com/google/uuid"
)

// *REQUESTS
// *START DATE IS THE OLDEST DATE AND END DATE IS THE NEWEST DATE
type Req_Get_Daily_Calorie_List struct {
	Start_Date time.Time `json:"start_date" validate:"required"`
	End_Date   time.Time `json:"end_date" validate:"required"`
}

type Res_Get_Daily_Calorie_List struct {
	ID           uint      `json:"id"`
	Account_Id   uuid.UUID `json:"account_id"`
	Calories     float32   `json:"calories"`
	Date_Created time.Time `json:"date_created"`
}
