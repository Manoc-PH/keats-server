package schemas

import "github.com/google/uuid"

type Req_Get_Recipe_Ingredients struct {
	Recipe_Id uuid.UUID `json:"recipe_id" validate:"required"`
}
type Res_Get_Recipe_Ingredients struct {
	Ingredients []Recipe_Ingredient_Details_Schema `json:"ingredients"`
}

type Recipe_Ingredient_Details_Schema struct {
	ID                    uuid.UUID `json:"id"`
	Ingredient_Mapping_Id uuid.UUID `json:"ingredient_mapping_id"`
	Food_Id               uuid.UUID `json:"food_id"`
	Calories              float32   `json:"calories"`
	Amount                float32   `json:"amount"`
	Amount_Unit           string    `json:"amount_unit"`
	Amount_Unit_Desc      string    `json:"amount_unit_desc"`
	Serving_Size          float32   `json:"serving_size"`
	Recipe_Id             uuid.UUID `json:"recipe_id"`
	Name                  string    `json:"name"`
	Name_Ph               string    `json:"name_ph"`
	Name_Owner            string    `json:"name_owner"`
}
