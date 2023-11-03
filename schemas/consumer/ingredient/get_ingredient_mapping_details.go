package schemas

import (
	"server/models"

	"github.com/google/uuid"
)

// *REQUESTS
type Req_Get_Ingredient_Mapping_Details struct {
	Ingredient_Mapping_ID uuid.UUID `json:"ingredient_mapping_id" validate:"required"`
	Return_Mappings       bool      `json:"return_mappings"`
}

type Res_Get_Ingredient_Mapping_Details struct {
	Ingredient            models.Ingredient            `json:"ingredient"`
	Ingredient_Variant    models.Ingredient_Variant    `json:"ingredient_variant"`
	Ingredient_Subvariant models.Ingredient_Subvariant `json:"ingredient_subvariant"`
	Nutrient              models.Nutrient              `json:"nutrient"`
	Ingredient_Images     []models.Ingredient_Image    `json:"ingredient_images"`
}
