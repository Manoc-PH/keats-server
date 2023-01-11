package models

type Recipe_Step struct {
	ID         int    `json:"id"`
	Recipe_Id  int    `json:"recipe_id"`
	Step_Order int    `json:"step_order"`
	Step_Desc  string `json:"step_desc"`
}
