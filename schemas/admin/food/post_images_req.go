package schemas

import (
	"server/models"
	"time"
)

type Req_Post_Images_Req struct {
	Food_Images []models.Food_Image `json:"food_images" validate:"required"`
	Timestamp   time.Time           `json:"timestamp" validate:"required"`
}
type Res_Post_Images_Req struct {
	Food_Images []models.Food_Image `json:"food_images"`
	Signature   string              `json:"signature"`
	Timestamp   string              `json:"timestamp"`
}
