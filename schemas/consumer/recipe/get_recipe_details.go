package schemas

import (
	"server/models"

	"github.com/google/uuid"
)

type Req_Get_Recipe_Details struct {
	Recipe_Id uuid.UUID `json:"recipe_id" validate:"required"`
}
type Res_Get_Recipe_Details struct {
	Recipe        models.Recipe         `json:"recipe"`
	Recipe_Images []models.Recipe_Image `json:"recipe_images"`
	Nutrients     models.Nutrient       `json:"nutrients"`
}
