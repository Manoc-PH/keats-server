package models

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID                 uint      `json:"id"`
	Owner_Id           uuid.UUID `json:"owner_id"`
	Recipe_Nutrient_Id uint      `json:"recipe_nutrient_id"`
	Name               string    `json:"name "`
	Name_Owner         string    `json:"name_owner"`
	Main_Image_Link    string    `json:"main_image_link"`
	Prep_Mins          int       `json:"prep_mins"`
	Servings           int       `json:"servings"`
	Saves              int       `json:"saves"`
	Date_Created       time.Time `json:"date_created"`
	Date_Updated       time.Time `json:"date_updated"`
	// time.Time SHOULD BE IN ISO STRING
}

type Recipe_Nutrient struct {
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

type Recipe_Ingredient struct {
	ID                          uint    `json:"id"`
	Food_Id                     uint    `json:"food_id"`
	Ingredient_Amount           float32 `json:"ingredient_amount"`
	Ingredient_Amount_Unit      string  `json:"ingredient_amount_unit"`
	Ingredient_Amount_Unit_Desc string  `json:"ingredient_amount_unit_desc"`
	Ingredient_Serving_Size     float32 `json:"ingredient_serving_size"`
}

type Recipe_Step struct {
	ID         uint   `json:"id"`
	Recipe_Id  uint   `json:"recipe_id"`
	Step_Order int    `json:"step_order"`
	Step_Desc  string `json:"step_desc"`
}
