package schemas

type Req_Get_Recipe_Ingredients struct {
	Recipe_Id uint `json:"recipe_id" validate:"required"`
}
type Res_Get_Recipe_Ingredients struct {
	Ingredients []Recipe_Ingredient_Details_Schema `json:"ingredients"`
}

type Recipe_Ingredient_Details_Schema struct {
	ID                    uint    `json:"id"`
	Ingredient_Mapping_Id uint    `json:"ingredient_mapping_id"`
	Food_Id               uint    `json:"food_id"`
	Amount                float32 `json:"amount"`
	Amount_Unit           string  `json:"amount_unit"`
	Amount_Unit_Desc      string  `json:"amount_unit_desc"`
	Serving_Size          float32 `json:"serving_size"`
	Recipe_Id             uint    `json:"recipe_id"`
	Name                  string  `json:"name"`
	Name_Ph               string  `json:"name_ph"`
	Name_Owner            string  `json:"name_owner"`
}
