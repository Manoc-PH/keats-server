package schemas

import (
	"server/models"
	"time"
)

// *START DATE IS THE OLDEST DATE AND END DATE IS THE NEWEST DATE
type Req_Get_Macros_List struct {
	Start_Date time.Time `json:"start_date" validate:"required"`
	End_Date   time.Time `json:"end_date" validate:"required"`
}

type Req_Post_Intake struct {
	Food_Id          uint    `json:"food_id" validate:"required_if=Recipe_Id 0"`
	Recipe_Id        uint    `json:"recipe_id" validate:"required_if=Food_Id 0"`
	Amount           float32 `json:"amount" validate:"required"`
	Amount_Unit      string  `json:"amount_unit" validate:"oneof='g' 'ml'"`
	Amount_Unit_Desc string  `json:"amount_unit_desc" validate:"oneof='gram' 'milliliter'"`
	Serving_Size     float32 `json:"serving_size"`
}

type Added_Macros struct {
	Calories int `json:"calories"`
	Protein  int `json:"protein"`
	Carbs    int `json:"carbs"`
	Fats     int `json:"fats"`
}
type Added_Coins_And_XP struct {
	Coins uint `json:"coins"`
	XP    uint `json:"xp"`
}
type Res_Post_Intake struct {
	Added_Macros       Added_Macros       `json:"added_macros"`
	Added_Coins_And_XP Added_Coins_And_XP `json:"added_coins_and_xp"`
	Intake             models.Intake      `json:"intake"`
	Food               models.Food        `json:"food"`
}
