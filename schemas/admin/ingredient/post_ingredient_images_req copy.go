package schemas

import (
	"server/models"
)

type Req_Post_Ingredient_Images_Confirm struct {
	Ingredient_Images []models.Ingredient_Image `json:"ingredient_images"`
}
