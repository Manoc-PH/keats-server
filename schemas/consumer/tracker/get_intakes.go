package schemas

import (
	"time"

	"github.com/google/uuid"
)

//* RESPONSE
type Res_Get_Intakes struct {
	ID                    uuid.UUID `json:"id"`
	Account_Id            uuid.UUID `json:"account_id"`
	Ingredient_Mapping_Id uuid.UUID `json:"ingredient_mapping_id"`
	Food_Id               uuid.UUID `json:"food_id"`
	Date_Created          time.Time `json:"date_created"`
	// Calories              float32   `json:"calories"`
	Amount           float32 `json:"amount"`
	Amount_Unit      string  `json:"amount_unit"`
	Amount_Unit_Desc string  `json:"amount_unit_desc"`
	Serving_Size     float32 `json:"serving_size"`

	// Ingredient
	Ingredient_Id                 uuid.UUID `json:"ingredient_id"`
	Ingredient_Name               string    `json:"ingredient_name"`
	Ingredient_Name_Ph            string    `json:"ingredient_name_ph"`
	Ingredient_Variant_Name       string    `json:"ingredient_variant_name"`
	Ingredient_Variant_Name_Ph    string    `json:"ingredient_variant_name_ph"`
	Ingredient_Subvariant_Name    string    `json:"ingredient_subvariant_name"`
	Ingredient_Subvariant_Name_Ph string    `json:"ingredient_subvariant_name_ph"`
	Ingredient_Name_Owner         string    `json:"ingredient_name_owner"`

	// Food
	Food_Name       string `json:"food_name"`
	Food_Name_Ph    string `json:"food_name_ph"`
	Food_Name_Owner string `json:"food_name_owner"`
}
