package schemas

import "server/models"

// *REQUESTS
type Req_Delete_Intake struct {
	Intake_ID uint `json:"intake_id" validate:"required"`
}

type Deleted_Daily_Nutrients struct {
	Calories      float32 `json:"calories"`
	Protein       float32 `json:"protein"`
	Carbs         float32 `json:"carbs"`
	Fats          float32 `json:"fats"`
	Trans_Fat     float32 `json:"trans_fat"`
	Saturated_Fat float32 `json:"saturated_fat"`
	Sugars        float32 `json:"sugars"`
	Sodium        float32 `json:"sodium"`
}
type Res_Delete_Intake struct {
	Deleted_Daily_Nutrients Deleted_Daily_Nutrients `json:"deleted_macros"`
	// Deleted_Coins_And_XP    Added_Coins_And_XP      `json:"deleted_coins_and_xp"`
	Intake models.Intake `json:"intake"`
	Food   models.Food   `json:"food"`
}
