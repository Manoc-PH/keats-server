package models

import (
	"time"

	"github.com/google/uuid"
)

type Macros struct {
	ID              int       `json:"id"`
	Account_Id      uuid.UUID `json:"account_id"`
	Date_Created    time.Time `json:"date_created"`
	Calories        float32   `json:"calories"`
	Protein         float32   `json:"protein"`
	Carbs           float32   `json:"carbs"`
	Fats            float32   `json:"fats"`
	Total_Calories  float32   `json:"total_calories"`
	Total_Protein   float32   `json:"total_protein"`
	Total_Carbs     float32   `json:"total_carbs"`
	Total_Fats      float32   `json:"total_fats"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id"`
	// time.Time SHOULD BE IN ISO STRING
}
