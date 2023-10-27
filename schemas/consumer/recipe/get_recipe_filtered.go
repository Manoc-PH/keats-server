package schemas

import "server/models"

type Req_Get_Recipe_Filtered struct {
	Filter  string `json:"filter" validate:"required,oneof=h_protein h_carbs h_fats l_cal l_carbs l_fats"`
	Created bool   `json:"created"`
	Liked   bool   `json:"liked"`
}

type Res_Get_Recipe_Filtered struct {
	Recipes []models.Recipe `json:"recipes"`
}
