package schemas

import (
	"server/models"
)

// *REQUESTS
type Req_Get_Food_Details struct {
	Food_ID uint `json:"food_id" validate:"required"`
}

type Res_Get_Food_Details struct {
	Food        models.Food         `json:"food"`
	Nutrient    models.Nutrient     `json:"nutrient"`
	Food_Images []models.Food_Image `json:"food_images"`
}
