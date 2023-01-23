package schemas

import "server/models"

// *REQUESTS
type Req_Delete_Intake struct {
	Intake_ID uint `json:"intake_id" validate:"required"`
}

type Deleted_Macros struct {
	Calories int `json:"calories"`
	Protein  int `json:"protein"`
	Carbs    int `json:"carbs"`
	Fats     int `json:"fats"`
}
type Res_Delete_Intake struct {
	Deleted_Macros       Deleted_Macros     `json:"deleted_macros"`
	Deleted_Coins_And_XP Added_Coins_And_XP `json:"deleted_coins_and_xp"`
	Intake               models.Intake      `json:"intake"`
	Food                 models.Food        `json:"food"`
}
