package models

import "time"

type Food struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Name_Ph          string    `json:"name_ph"`
	Name_Brand       string    `json:"name_brand"`
	Food_Nutrient_Id int       `json:"food_nutrient_id"`
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
}
