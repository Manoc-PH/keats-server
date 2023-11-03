package schemas

import "github.com/google/uuid"

// *REQUESTS
type Req_Patch_Recipe struct {
	Recipe              Recipe_Patch               `json:"recipe" validate:"required,dive"`
	Recipe_Ingredients  []Recipe_Patch_Ingredient  `json:"recipe_ingredients" validate:"required,max=10,dive"`
	Recipe_Instructions []Recipe_Patch_Instruction `json:"recipe_instructions" validate:"required,max=20,dive"`
}

// *RESPONSES
type Res_Patch_Recipe struct {
	Recipe Recipe_Patch `json:"recipe"`
}

// Schemas
type Recipe_Patch struct {
	ID             uuid.UUID `json:"id" validate:"required"`
	Name           string    `json:"name" validate:"required"`
	Name_Ph        string    `json:"name_ph"`
	Category_Id    uuid.UUID `json:"category_id"`
	Thumbnail_URL  string    `json:"thumbnail_url"`
	Thumbnail_Name string    `json:"thumbnail_name"`
	Image_URL      string    `json:"image_url"`
	Image_Name     string    `json:"image_name"`
	Servings       uint      `json:"servings" validate:"required"`
	Servings_Size  float32   `json:"servings_size" validate:"required"`
	Prep_Time      uint      `json:"prep_time" validate:"required"`
	Description    string    `json:"description"`
	Nutrient_Id    uuid.UUID `json:"nutrient_id"`
}

// Client can send either delete, update, or insert action type which will determine what to do with the data
type Recipe_Patch_Ingredient struct {
	ID                    uuid.UUID `json:"id" validate:"required_if=Action_Type delete,required_if=Action_Type update"`
	Food_Id               uuid.UUID `json:"food_id" validate:"required_if=Ingredient_Mapping_Id 0"`
	Ingredient_Mapping_Id uuid.UUID `json:"ingredient_mapping_id" validate:"required_if=Food_Id 0"`
	Amount                float32   `json:"amount" validate:"required"`
	Amount_Unit           string    `json:"amount_unit" validate:"required"`
	Amount_Unit_Desc      string    `json:"amount_unit_desc" validate:"required"`
	Serving_Size          float32   `json:"serving_size"`
	Recipe_Id             uuid.UUID `json:"recipe_id"`
	Action_Type           string    `json:"action_type"  validate:"required,oneof=delete update insert"`
}

type Recipe_Patch_Instruction struct {
	ID                      uuid.UUID `json:"id" validate:"required_if=Action_Type delete,required_if=Action_Type update"`
	Recipe_Id               uuid.UUID `json:"recipe_id"`
	Instruction_Description string    `json:"instruction_description" validate:"required"`
	Step_Num                uint      `json:"step_num"`
	Action_Type             string    `json:"action_type" validate:"required,oneof=delete update insert"`
}
