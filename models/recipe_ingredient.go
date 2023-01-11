package models

type Recipe_Ingredient struct {
	ID                          int     `json:"id"`
	Recipe_Id                   int     `json:"recipe_id"`
	Food_Id                     int     `json:"food_id"`
	Ingredient_Amount           float32 `json:"ingredient_amount"`
	Ingredient_Amount_Unit      string  `json:"ingredient_amount_unit"`
	Ingredient_Amount_Unit_Desc string  `json:"ingredient_amount_unit_desc"`
	Ingredient_Serving_Size     float32 `json:"ingredient_serving_size"`
}
