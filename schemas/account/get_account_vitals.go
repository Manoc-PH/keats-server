package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Res_Get_Account_Vitals struct {
	Account_ID         uuid.UUID `json:"account_id"`
	Weight             uint      `json:"weight"`
	Height             uint      `json:"height"`
	Birthday           time.Time `json:"birthday"`
	Age                uint      `json:"age"`
	Sex                string    `json:"sex"`
	Activity_Lvl_Id    uuid.UUID `json:"activity_lvl_id"`
	Activity_Lvl_Name  string    `json:"activity_lvl_name"`
	Bmr_Multiplier     string    `json:"bmr_multiplier"`
	Diet_Plan_Id       uuid.UUID `json:"diet_plan_id"`
	Diet_Plan_Name     string    `json:"diet_plan_name"`
	Calorie_Percentage int       `json:"calorie_percentage"`
	Protein_Percentage int       `json:"protein_percentage"`
	Fats_Percentage    int       `json:"fats_percentage"`
	Carbs_Percentage   int       `json:"carbs_percentage"`
}
