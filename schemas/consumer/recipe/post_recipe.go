package schemas

import (
	"time"

	"github.com/google/uuid"
)

// *REQUESTS
type Req_Post_Recipe struct {
	Recipe              Recipe_Schema               `json:"recipe"`
	Recipe_Ingredients  []Recipe_Ingredient_Schema  `json:"recipe_ingredients" validate:"required"`
	Recipe_Instructions []Recipe_Instruction_Schema `json:"recipe_instructions" validate:"required"`
	Timestamp           time.Time                   `json:"timestamp" validate:"required"`
}

// *RESPONSES
type Res_Post_Recipe struct {
	Recipe              Recipe_Schema               `json:"recipe"`
	Recipe_Ingredients  []Recipe_Ingredient_Schema  `json:"recipe_ingredients" validate:"required"`
	Recipe_Instructions []Recipe_Instruction_Schema `json:"recipe_instructions" validate:"required"`
	Signature           string                      `json:"signature"`
	Timestamp           string                      `json:"timestamp"`
}

// Schemas
type Recipe_Schema struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name" validate:"required"`
	Name_Ph              string    `json:"name_ph"`
	Name_Owner           string    `json:"name_owner" validate:"required"`
	Owner_Id             uuid.UUID `json:"owner_id" validate:"required"`
	Date_Created         time.Time `json:"date_created"`
	Category_Id          uint      `json:"category_id"`
	Thumbnail_Image_Link string    `json:"thumbnail_image_link"`
	Main_Image_Link      string    `json:"main_image_link"`
	Likes                uint      `json:"likes"`
	Rating               uint      `json:"rating"`
	Servings             uint      `json:"servings" validate:"required"`
	Servings_Size        uint      `json:"servings_size" validate:"required"`
	Prep_Time            uint      `json:"prep_time" validate:"required"`
	Description          string    `json:"description"`
}
type Recipe_Ingredient_Schema struct {
	ID                    uint    `json:"id"`
	Food_Id               uint    `json:"food_id" validate:"required_if=Ingredient_Mapping_Id 0"`
	Ingredient_Mapping_Id uint    `json:"ingredient_mapping_id" validate:"required_if=Food_Id 0"`
	Amount                float32 `json:"amount" validate:"required"`
	Amount_Unit           string  `json:"amount_unit"`
	Amount_Unit_Desc      string  `json:"amount_unit_desc"`
	Serving_Size          float32 `json:"serving_size"`
}
type Recipe_Instruction_Schema struct {
	ID                      uint   `json:"id"`
	Recipe_Id               uint   `json:"recipe_id"`
	Instruction_Description string `json:"instruction_description" validate:"required"`
	Step_Num                uint   `json:"step_num"`
}
