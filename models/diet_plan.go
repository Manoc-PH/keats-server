package models

import (
	"github.com/google/uuid"
)

type Diet_Plan struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Main_Image_Link    string    `json:"main_image_link"`
	Background_Color   string    `json:"background_color"`
	Diet_Plan_Desc     string    `json:"diet_plan_desc"`
	Calorie_Percentage int       `json:"calorie_percentage"`
	Protein_Percentage int       `json:"protein_percentage"`
	Fats_Percentage    int       `json:"fats_percentage"`
	Carbs_Percentage   int       `json:"carbs_percentage"`
}
