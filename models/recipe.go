package models

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID              uint      `json:"id"`
	Owner_Id        uuid.UUID `json:"owner_id"`
	Main_Image_Link string    `json:"main_image_link"`
	Prep_Mins       int       `json:"prep_mins"`
	Servings        int       `json:"servings"`
	Saves           int       `json:"saves"`
	Date_Created    time.Time `json:"date_created"`
	Date_Updated    time.Time `json:"date_updated"`
	// time.Time SHOULD BE IN ISO STRING
}

type Recipe_Step struct {
	ID         uint   `json:"id"`
	Recipe_Id  uint   `json:"recipe_id"`
	Step_Order int    `json:"step_order"`
	Step_Desc  string `json:"step_desc"`
}

type Recipe_Ingredient struct {
	ID                          uint    `json:"id"`
	Recipe_Id                   uint    `json:"recipe_id"`
	Food_Id                     uint    `json:"food_id"`
	Ingredient_Amount           float32 `json:"ingredient_amount"`
	Ingredient_Amount_Unit      string  `json:"ingredient_amount_unit"`
	Ingredient_Amount_Unit_Desc string  `json:"ingredient_amount_unit_desc"`
	Ingredient_Serving_Size     float32 `json:"ingredient_serving_size"`
}
