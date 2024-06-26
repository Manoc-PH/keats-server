package schemas

import (
	"server/models"

	"github.com/google/uuid"
)

// *REQUESTS
type Req_Post_Intake struct {
	Food_Id               uuid.UUID `json:"food_id" validate:"required_if=Ingredient_Mapping_Id 0"`
	Ingredient_Mapping_Id uuid.UUID `json:"ingredient_mapping_id" validate:"required_if=Food_Id 0"`
	Amount                float32   `json:"amount" validate:"required"`
	Amount_Unit           string    `json:"amount_unit" validate:"oneof='g' 'ml'"`
	Amount_Unit_Desc      string    `json:"amount_unit_desc" validate:"oneof='gram' 'milliliter'"`
	Serving_Size          float32   `json:"serving_size"`
}

// type Added_Coins_And_XP struct {
// 	Coins int `json:"coins"`
// 	XP    int `json:"xp"`
// }
type Ingredient_Mapping_Schema struct {
	Ingredient            models.Ingredient            `json:"ingredient"`
	Ingredient_Variant    models.Ingredient_Variant    `json:"ingredient_variant"`
	Ingredient_Subvariant models.Ingredient_Subvariant `json:"ingredient_subvariant"`
	// TODO Check if nutrient is needed, if not remove
	Nutrient models.Nutrient `json:"nutrient"`
}
type Food_Mapping_Schema struct {
	Food     models.Food     `json:"food"`
	Nutrient models.Nutrient `json:"nutrient"`
}
type Res_Post_Intake struct {
	Added_Daily_Nutrients models.Nutrient `json:"added_daily_nutrients"`
	// Added_Coins_And_XP    Added_Coins_And_XP    `json:"added_coins_and_xp"`
	Intake     models.Intake              `json:"intake"`
	Ingredient *Ingredient_Mapping_Schema `json:"ingredient"`
	Food       *models.Food               `json:"food"`
}
