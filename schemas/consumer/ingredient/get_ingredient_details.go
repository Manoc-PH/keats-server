package schemas

import (
	"server/models"
)

// *REQUESTS
type Req_Get_Ingredient_Details struct {
	Ingredient_ID uint `json:"ingredient_id" validate:"required"`
}
type Ingredient_Mapping_With_Name struct {
	ID                         uint   `json:"id"`
	Ingredient_Id              uint   `json:"ingredient_id"`
	Ingredient_Variant_Id      uint   `json:"ingredient_variant_id"`
	Ingredient_Subvariant_Id   uint   `json:"ingredient_subvariant_id"`
	Nutrient_Id                uint   `json:"nutrient_id"`
	Ingredient_Variant_Name    string `json:"ingredient_variant_name"`
	Ingredient_Subvariant_Name string `json:"ingredient_subvariant_name"`
}
type Res_Get_Ingredient_Details struct {
	Ingredient            models.Ingredient              `json:"ingredient"`
	Ingredient_Variant    models.Ingredient_Variant      `json:"ingredient_variant"`
	Ingredient_Subvariant models.Ingredient_Subvariant   `json:"ingredient_subvariant"`
	Nutrient              models.Nutrient                `json:"nutrient"`
	Ingredient_Mapping_ID uint                           `json:"ingredient_mapping_id"`
	Ingredient_Images     []models.Ingredient_Image      `json:"ingredient_images"`
	Ingredient_Mappings   []Ingredient_Mapping_With_Name `json:"ingredient_mappings"`
}
