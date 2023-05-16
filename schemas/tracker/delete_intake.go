package schemas

import "server/models"

// *REQUESTS
type Req_Delete_Intake struct {
	Intake_ID uint `json:"intake_id" validate:"required"`
}

type Res_Delete_Intake struct {
	Deleted_Daily_Nutrients models.Nutrient `json:"deleted_daily_nutrients"`
	Intake                  models.Intake   `json:"intake"`
}
