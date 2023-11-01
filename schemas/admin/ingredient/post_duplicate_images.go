package schemas

import (
	"server/models"

	"github.com/google/uuid"
)

type Req_Post_Duplicate_Images struct {
	Ingredient_Mapping_Id        uuid.UUID `json:"ingredient_mapping_id" validate:"required"`
	Copied_Ingredient_Mapping_Id uuid.UUID `json:"copied_ingredient_mapping_id" validate:"required"`
}
type Res_Post_Duplicate_Images struct {
	Ingredient_Images []models.Ingredient_Image `json:"ingredient_images"`
}
