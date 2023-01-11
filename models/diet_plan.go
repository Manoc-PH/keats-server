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
	Calorie_Percentage int16     `json:"calorie_percentage"`
	Protein_Percentage int16     `json:"protein_percentage"`
	Fats_Percentage    int16     `json:"fats_percentage"`
	Carbs_Percentage   int16     `json:"carbs_percentage"`
}
