package schemas

import "server/models"

type Ingredient_Details struct {
	Ingredient            models.Ingredient            `json:"ingredient"`
	Ingredient_Variant    models.Ingredient_Variant    `json:"ingredient_variant"`
	Ingredient_Subvariant models.Ingredient_Subvariant `json:"ingredient_subvariant"`
	Nutrient              models.Nutrient              `json:"nutrient"`
	Ingredient_Images     []models.Ingredient_Image    `json:"ingredient_images"`
}

type Req_Put_Ingredient_Details struct {
	Ingredient_Mapping_ID uint               `json:"Ingredient_Mapping_ID" validate:"requried"`
	Ingredient_Details    Ingredient_Details `json:"ingredient_details"`
}
