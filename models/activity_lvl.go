package models

import (
	"github.com/google/uuid"
)

type Activity_Lvl struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Main_Image_Link   string    `json:"main_image_link"`
	Background_Color  string    `json:"background_color"`
	Activity_Lvl_Desc string    `json:"activity_lvl_desc"`
	Bmr_Multipler     float32   `json:"bmr_multipler"`
}
