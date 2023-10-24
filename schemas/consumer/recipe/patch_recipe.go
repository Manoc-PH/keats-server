package schemas

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
	ID                   uint    `json:"id" validate:"required"`
	Name                 string  `json:"name" validate:"required"`
	Name_Ph              string  `json:"name_ph"`
	Category_Id          uint    `json:"category_id"`
	Thumbnail_Image_Link string  `json:"thumbnail_image_link"`
	Main_Image_Link      string  `json:"main_image_link"`
	Servings             uint    `json:"servings" validate:"required"`
	Servings_Size        float32 `json:"servings_size" validate:"required"`
	Prep_Time            uint    `json:"prep_time" validate:"required"`
	Description          string  `json:"description"`
}

// Client can send either delete, update, or insert action type which will determine what to do with the data
type Recipe_Patch_Ingredient struct {
	ID                    uint    `json:"id" validate:"required_if=Action_Type delete,required_if=Action_Type update"`
	Food_Id               uint    `json:"food_id" validate:"required_if=Ingredient_Mapping_Id 0"`
	Ingredient_Mapping_Id uint    `json:"ingredient_mapping_id" validate:"required_if=Food_Id 0"`
	Amount                float32 `json:"amount" validate:"required"`
	Amount_Unit           string  `json:"amount_unit"`
	Amount_Unit_Desc      string  `json:"amount_unit_desc"`
	Serving_Size          float32 `json:"serving_size"`
	Recipe_Id             uint    `json:"recipe_id"`
	Action_Type           string  `json:"action_type"  validate:"required,oneof=delete update insert"`
}

type Recipe_Patch_Instruction struct {
	ID                      uint   `json:"id" validate:"required_if=Action_Type delete,required_if=Action_Type update"`
	Recipe_Id               uint   `json:"recipe_id"`
	Instruction_Description string `json:"instruction_description" validate:"required"`
	Step_Num                uint   `json:"step_num"`
	Action_Type             string `json:"action_type" validate:"required,oneof=delete update insert"`
}
