package schemas

import "github.com/google/uuid"

type Ingredient_Image_Confirm struct {
	ID                    uuid.UUID `json:"id" validate:"required"`
	Ingredient_Mapping_Id uuid.UUID `json:"ingredient_mapping_id"`
	Name_File             string    `json:"name_file" validate:"required"`
	Name_URL              string    `json:"name_url" validate:"required"`
	Amount                float32   `json:"amount"`
	Amount_Unit           string    `json:"amount_unit"`
	Amount_Unit_Desc      string    `json:"amount_unit_desc"`
}

type Req_Post_Images_Confirm struct {
	Ingredient_Images []Ingredient_Image_Confirm `json:"ingredient_images" validate:"required"`
}
