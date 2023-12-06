package schemas

import "server/models"

type Req_Get_Recipe_Filtered struct {
	Filter  string `json:"filter" validate:"required_if_empty=Created Liked,oneof=h_protein h_carbs h_fats l_cal l_carbs l_fats"`
	Created bool   `json:"created" validate:"required_if_empty=Filter Liked"`
	Liked   bool   `json:"liked" validate:"required_if_empty=Filter Created"`
}

type Res_Get_Recipe_Filtered struct {
	Recipes []models.Recipe `json:"recipes"`
}
