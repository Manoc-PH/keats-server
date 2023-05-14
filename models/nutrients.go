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
	Trans_Fat       float32   `json:"trans_fat"`
	Saturated_Fat   float32   `json:"saturated_fat"`
	Sugars          float32   `json:"sugars"`
	Fiber           float32   `json:"fiber"`
	Sodium          float32   `json:"sodium"`
	Iron            float32   `json:"iron"`
	Calcium         float32   `json:"calcium"`
	// time.Time SHOULD BE IN ISO STRING
}

type Nutrient struct {
	ID                    uint    `json:"id"`
	Food_Id               uint    `json:"food_id"`
	Ingredient_Mapping_Id uint    `json:"ingredient_mapping_id"`
	Amount                float32 `json:"amount"`
	Amount_Unit           string  `json:"amount_unit"`
	Amount_Unit_Desc      string  `json:"amount_unit_desc"`
	Serving_Size          float32 `json:"serving_size"`
	Calories              float32 `json:"calories"`
	Protein               float32 `json:"protein"`
	Carbs                 float32 `json:"carbs"`
	Fats                  float32 `json:"fats"`
	//
	Trans_Fat     float32 `json:"trans_fat"`
	Saturated_Fat float32 `json:"saturated_fat"`
	Sugars        float32 `json:"sugars"`
	Fiber         float32 `json:"fiber"`
	Sodium        float32 `json:"sodium"`
	Iron          float32 `json:"iron"`
	Calcium       float32 `json:"calcium"`
}
