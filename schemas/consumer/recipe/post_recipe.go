package schemas

import (
	"server/models"
	"time"

	"github.com/google/uuid"
)

// *REQUESTS
type Req_Post_Recipe struct {
	Recipe              Recipe_Post               `json:"recipe"`
	Recipe_Ingredients  []Recipe_Ingredient_Post  `json:"recipe_ingredients" validate:"required,max=10,dive"`
	Recipe_Instructions []Recipe_Instruction_Post `json:"recipe_instructions" validate:"required,max=20,dive"`
}

// *RESPONSES
type Res_Post_Recipe struct {
	Recipe              Recipe_Post               `json:"recipe"`
	Recipe_Ingredients  []Recipe_Ingredient_Post  `json:"recipe_ingredients"`
	Recipe_Instructions []Recipe_Instruction_Post `json:"recipe_instructions"`
	Nutrient            models.Nutrient           `json:"nutrient"`
}

// Schemas
type Recipe_Post struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name" validate:"required"`
	Name_Ph        string    `json:"name_ph"`
	Name_Owner     string    `json:"name_owner"`
	Owner_Id       uuid.UUID `json:"owner_id"`
	Date_Created   time.Time `json:"date_created"`
	Category_Id    uuid.UUID `json:"category_id"`
	Nutrient_Id    uuid.UUID `json:"nutrient_id"`
	Thumbnail_URL  string    `json:"thumbnail_url"`
	Thumbnail_Name string    `json:"thumbnail_name"`
	Image_URL      string    `json:"image_url"`
	Image_Name     string    `json:"image_name"`
	Likes          uint      `json:"likes"`
	Rating         float32   `json:"rating"`
	Rating_Count   uint      `json:"rating_count"`
	Servings       float32   `json:"servings" validate:"required"`
	Servings_Size  float32   `json:"servings_size"`
	Prep_Time      float32   `json:"prep_time" validate:"required"`
	Description    string    `json:"description"`
}
type Recipe_Ingredient_Post struct {
	ID                    uuid.UUID `json:"id"`
	Food_Id               uuid.UUID `json:"food_id" validate:"required_if=Ingredient_Mapping_Id 0"`
	Ingredient_Mapping_Id uuid.UUID `json:"ingredient_mapping_id" validate:"required_if=Food_Id 0"`
	Amount                float32   `json:"amount" validate:"required"`
	Amount_Unit           string    `json:"amount_unit"`
	Amount_Unit_Desc      string    `json:"amount_unit_desc"`
	Serving_Size          float32   `json:"serving_size"`
	Recipe_Id             uuid.UUID `json:"recipe_id"`
}
type Recipe_Instruction_Post struct {
	ID                      uuid.UUID `json:"id"`
	Recipe_Id               uuid.UUID `json:"recipe_id"`
	Instruction_Description string    `json:"instruction_description" validate:"required"`
	Step_Num                uint      `json:"step_num" validate:"required"`
}

// *SAMPLE RESPONSE
// "recipe": {
// 	"id": 6,
// 	"name": "test",
// 	"name_ph": "",
// 	"name_owner": "Cloyd Abad",
// 	"owner_id": "4767bca7-4911-4496-9de2-fb6b2d318c6c",
// 	"date_created": "0001-01-01T00:00:00Z",
// 	"category_id": 0,
// 	"thumbnail_image_link": "",
// 	"main_image_link": "",
// 	"likes": 0,
// 	"rating": 0,
// 	"servings": 4,
// 	"servings_size": 40,
// 	"prep_time": 20,
// 	"description": ""
// },
// "recipe_ingredients": [
// 	{
// 			"id": 3,
// 			"food_id": 0,
// 			"ingredient_mapping_id": 89,
// 			"amount": 100,
// 			"amount_unit": "",
// 			"amount_unit_desc": "",
// 			"serving_size": 0
// 	},
// 	{
// 			"id": 4,
// 			"food_id": 0,
// 			"ingredient_mapping_id": 99,
// 			"amount": 100,
// 			"amount_unit": "",
// 			"amount_unit_desc": "",
// 			"serving_size": 0
// 	},
// 	{
// 			"id": 5,
// 			"food_id": 0,
// 			"ingredient_mapping_id": 29,
// 			"amount": 100,
// 			"amount_unit": "",
// 			"amount_unit_desc": "",
// 			"serving_size": 0
// 	}
// ],
// "recipe_instructions": [
// 	{
// 			"id": 2,
// 			"recipe_id": 0,
// 			"instruction_description": "test",
// 			"step_num": 0
// 	},
// 	{
// 			"id": 3,
// 			"recipe_id": 0,
// 			"instruction_description": "test2",
// 			"step_num": 0
// 	},
// 	{
// 			"id": 4,
// 			"recipe_id": 0,
// 			"instruction_description": "test3",
// 			"step_num": 0
// 	}
// ],
// "signature": "f3d52cedcab2c2383cc5ba6461b21525532bb047",
// "timestamp": "1697470437"
