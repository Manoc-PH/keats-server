package schemas

import "time"

// *START DATE IS THE OLDEST DATE AND END DATE IS THE NEWEST DATE
type Req_Get_Macros_List struct {
	Start_Date time.Time `json:"start_date" validate:"required"`
	End_Date   time.Time `json:"end_date" validate:"required"`
}

type Req_Post_Intake struct {
	Food_Id          uint    `json:"food_id" validate:"required_if=Recipe_Id nil"`
	Recipe_Id        uint    `json:"recipe_id" validate:"required_if=Food_Id nil"`
	Amount           float32 `json:"amount" validate:"required"`
	Amount_Unit      string  `json:"amount_unit" validate:"required"`
	Amount_Unit_Desc string  `json:"amount_unit_desc" validate:"required"`
	Serving_Size     float32 `json:"serving_size"`
}
