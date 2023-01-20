package models

import "time"

type Food struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Name_Ph          string    `json:"name_ph"`
	Name_Brand       string    `json:"name_brand"`
	Food_Nutrient_Id uint      `json:"food_nutrient_id"`
	Date_Created     time.Time `json:"date_created"`
	// time.Time SHOULD BE IN ISO STRING
}

type Food_Nutrient struct {
	ID               uint    `json:"id"`
	Food_Id          uint    `json:"food_id"`
	Amount           float32 `json:"amount"`
	Amount_Unit      string  `json:"amount_unit"`
	Amount_Unit_Desc string  `json:"amount_unit_desc"`
	Serving_Size     float32 `json:"serving_size"`
	Calories         float32 `json:"calories"`
	Protein          float32 `json:"protein"`
	Carbs            float32 `json:"carbs"`
	Fats             float32 `json:"fats"`
	//
	Trans_Fat     float32 `json:"trans_fat"`
	Saturated_Fat float32 `json:"saturated_fat"`
	Sugars        float32 `json:"sugars"`
	Sodium        float32 `json:"sodium"`
}
type Food_Image struct {
	ID               uint    `json:"id"`
	Food_Id          uint    `json:"food_id"`
	Name_File        string  `json:"name_file"`
	Amount           float32 `json:"amount"`
	Amount_Unit      string  `json:"amount_unit"`
	Amount_Unit_Desc string  `json:"amount_unit_desc"`
}
