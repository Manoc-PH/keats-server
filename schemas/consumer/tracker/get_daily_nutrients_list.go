package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Res_Get_Daily_Nutrients_List struct {
	ID           uint      `json:"id"`
	Account_Id   uuid.UUID `json:"account_id"`
	Calories     float32   `json:"calories"`
	Date_Created time.Time `json:"date_created"`
}
