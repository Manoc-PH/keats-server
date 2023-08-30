package schemas

import (
	"server/models"
	"time"
)

type Req_Post_Ingredient_Images struct {
	Ingredient_Images []models.Ingredient_Image `json:"ingredient_images"`
	Timestamp         time.Time                 `json:"timestamp" validate:"required"`
}
type Res_Post_Ingredient_Images struct {
	Ingredient_Images []models.Ingredient_Image `json:"ingredient_images"`
	Signature         string                    `json:"signature"`
	Timestamp         string                    `json:"timestamp"`
}
