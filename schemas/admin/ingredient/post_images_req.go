package schemas

import (
	"time"

	"github.com/google/uuid"
)

// TODO ADD REQUIRED TO SOME FIELDS IN INGREDIENT IMAGE
type Ingredient_Image_Req struct {
	ID                    uuid.UUID `json:"id"`
	Ingredient_Mapping_Id uuid.UUID `json:"ingredient_mapping_id" validate:"required"`
	Name_File             string    `json:"name_file"`
	Name_URL              string    `json:"name_url"`
	Amount                float32   `json:"amount" validate:"required,gte=0"`
	Amount_Unit           string    `json:"amount_unit" validate:"required"`
	Amount_Unit_Desc      string    `json:"amount_unit_desc" validate:"required"`
}

type Req_Post_Images_Req struct {
	Ingredient_Images []Ingredient_Image_Req `json:"ingredient_images" validate:"required"`
	Timestamp         time.Time              `json:"timestamp" validate:"required"`
}
type Res_Post_Images_Req struct {
	Ingredient_Images []Ingredient_Image_Req `json:"ingredient_images"`
	Signature         string                 `json:"signature"`
	Timestamp         string                 `json:"timestamp"`
}
