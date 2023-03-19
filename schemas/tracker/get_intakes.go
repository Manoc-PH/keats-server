package schemas

import (
	"time"

	"github.com/google/uuid"
)

//* RESPONSE
type Res_Get_Intakes struct {
	ID               uint      `json:"id"`
	Account_Id       uuid.UUID `json:"account_id"`
	Food_Id          uint      `json:"food_id"`
	Recipe_Id        uint      `json:"recipe_id"`
	Date_Created     time.Time `json:"date_created"`
	Amount           float32   `json:"amount"`
	Amount_Unit      string    `json:"amount_unit"`
	Amount_Unit_Desc string    `json:"amount_unit_desc"`
	Serving_Size     float32   `json:"serving_size"`
	// Food
	Food_Name                 string  `json:"food_name"`
	Food_Name_Ph              string  `json:"food_name_ph"`
	Food_Name_Brand           string  `json:"food_name_brand"`
	Food_Nutrient_Id          uint    `json:"food_nutrient_id"`
	Food_Nutrient_Calories    float32 `json:"food_nutrient_calories"`
	Food_Nutrient_Amount      float32 `json:"food_nutrient_amount"`
	Food_Nutrient_Amount_Unit string  `json:"food_nutrient_amount_unit"`

	// Recipe
	Recipe_Name       string `json:"recipe_name"`
	Recipe_Name_Owner string `json:"recipe_name_owner"`
}
