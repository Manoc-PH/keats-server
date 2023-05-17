package schemas

import (
	"server/models"
)

// *REQUESTS
type Req_Get_Ingredient_Details struct {
	Ingredient_Mapping_ID uint `json:"ingredient_mapping_id" validate:"required"`
}

type Res_Get_Ingredient_Details struct {
	Ingredient            models.Ingredient            `json:"ingredient"`
	Ingredient_Variant    models.Ingredient_Variant    `json:"ingredient_variant"`
	Ingredient_Subvariant models.Ingredient_Subvariant `json:"ingredient_subvariant"`
	Nutrient              models.Nutrient              `json:"nutrient"`
	Ingredient_Images     []models.Ingredient_Image    `json:"ingredient_images"`
}
