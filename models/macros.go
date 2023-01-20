package models

import (
	"time"

	"github.com/google/uuid"
)

type Macros struct {
	ID              uint      `json:"id"`
	Account_Id      uuid.UUID `json:"account_id"`
	Date_Created    time.Time `json:"date_created"`
	Calories        int       `json:"calories"`
	Protein         int       `json:"protein"`
	Carbs           int       `json:"carbs"`
	Fats            int       `json:"fats"`
	Max_Calories    int       `json:"max_calories"`
	Max_Protein     int       `json:"max_protein"`
	Max_Carbs       int       `json:"max_carbs"`
	Max_Fats        int       `json:"max_fats"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id"`
	// time.Time SHOULD BE IN ISO STRING
}
