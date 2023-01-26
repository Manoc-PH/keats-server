package schemas

import "server/models"

// *REQUESTS
type Req_Put_Intake struct {
	Intake_ID        uint    `json:"intake_id" validate:"required"`
	Food_Id          uint    `json:"food_id" validate:"required_if=Recipe_Id 0"`
	Recipe_Id        uint    `json:"recipe_id" validate:"required_if=Food_Id 0"`
	Amount           float32 `json:"amount" validate:"required"`
	Amount_Unit      string  `json:"amount_unit" validate:"oneof='g' 'ml'"`
	Amount_Unit_Desc string  `json:"amount_unit_desc" validate:"oneof='gram' 'milliliter'"`
	Serving_Size     float32 `json:"serving_size"`
}

type Res_Put_Intake struct {
	Added_Macros       Added_Macros       `json:"added_macros"`
	Added_Coins_And_XP Added_Coins_And_XP `json:"added_coins_and_xp"`
	Intake             models.Intake      `json:"intake"`
	Food               models.Food        `json:"food"`
}