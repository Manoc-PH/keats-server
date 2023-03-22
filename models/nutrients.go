package models

import (
	"time"

	"github.com/google/uuid"
)

type Daily_Nutrients struct {
	ID              uint      `json:"id"`
	Account_Id      uuid.UUID `json:"account_id"`
	Date_Created    time.Time `json:"date_created"`
	Calories        float32   `json:"calories"`
	Protein         float32   `json:"protein"`
	Carbs           float32   `json:"carbs"`
	Fats            float32   `json:"fats"`
	Max_Calories    float32   `json:"max_calories"`
	Max_Protein     float32   `json:"max_protein"`
	Max_Carbs       float32   `json:"max_carbs"`
	Max_Fats        float32   `json:"max_fats"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id"`
	// time.Time SHOULD BE IN ISO STRING
}
