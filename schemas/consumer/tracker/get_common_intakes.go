package schemas

import (
	"time"
)

// *REQUESTS
// *START DATE IS THE OLDEST DATE AND END DATE IS THE NEWEST DATE
type Req_Get_Common_Intakes struct {
	Start_Date time.Time `json:"start_date" validate:"required"`
	End_Date   time.Time `json:"end_date" validate:"required"`
}
type Res_Get_Common_Intakes struct {
	Intakes []Intake_Details
}

type Intake_Details struct {
	Food_ID                       uint   `json:"food_id"`
	Ingredient_Mapping_ID         uint   `json:"ingredient_mapping_id"`
	Ingredient_Count              uint   `json:"ingredient_count"`
	Food_Count                    uint   `json:"food_count"`
	Ingredient_ID                 uint   `json:"ingredient_id"`
	Ingredient_Name               string `json:"ingredient_name"`
	Ingredient_Name_Ph            string `json:"ingredient_name_ph"`
	Ingredient_Name_Owner         string `json:"ingredient_name_owner"`
	Ingredient_Variant_ID         uint   `json:"ingredient_variant_id"`
	Ingredient_Variant_Name       string `json:"ingredient_variant_name"`
	Ingredient_Variant_Name_Ph    string `json:"ingredient_variant_name_ph"`
	Ingredient_Subvariant_ID      uint   `json:"ingredient_subvariant_id"`
	Ingredient_Subvariant_Name    string `json:"ingredient_subvariant_name"`
	Ingredient_Subvariant_Name_Ph string `json:"ingredient_subvariant_name_ph"`
	Food_Name                     string `json:"food_name"`
	Food_Name_Ph                  string `json:"food_name_ph"`
	Food_Name_Owner               string `json:"food_name_owner"`
}
