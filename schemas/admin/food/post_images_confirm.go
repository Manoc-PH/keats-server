package schemas

import (
	"server/models"
)

type Req_Post_Images_Confirm struct {
	Food_Images []models.Food_Image `json:"food_images" validate:"required"`
}
