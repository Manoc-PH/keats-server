package models

import (
	"time"

	"github.com/google/uuid"
)

type Macros struct {
	ID              int       `json:"id"`
	Account_Id      uuid.UUID `json:"account_id"`
	Date_Created    time.Time `json:"date_created"`
	Calories        int       `json:"calories"`
	Protein         int       `json:"protein"`
	Carbs           int       `json:"carbs"`
	Fats            int       `json:"fats"`
	Total_Calories  int       `json:"total_calories"`
	Total_Protein   int       `json:"total_protein"`
	Total_Carbs     int       `json:"total_carbs"`
	Total_Fats      int       `json:"total_fats"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id"`
	// time.Time SHOULD BE IN ISO STRING
}
