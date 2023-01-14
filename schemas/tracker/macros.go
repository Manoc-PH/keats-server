package schemas

import "time"

// *START DATE IS THE OLDEST DATE AND END DATE IS THE NEWEST DATE
type Req_Get_Macros_List struct {
	Start_Date time.Time `json:"start_date" validate:"required"`
	End_Date   time.Time `json:"end_date" validate:"required"`
}
