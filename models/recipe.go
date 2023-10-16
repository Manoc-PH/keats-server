package models

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Name_Ph              string    `json:"name_ph"`
	Name_Owner           string    `json:"name_owner"`
	Owner_Id             uuid.UUID `json:"owner_id"`
	Date_Created         time.Time `json:"date_created"`
	Category_Id          uint      `json:"category_id"`
	Thumbnail_Image_Link string    `json:"thumbnail_image_link"`
	Main_Image_Link      string    `json:"main_image_link"`
	Likes                uint      `json:"likes"`
	Rating               uint      `json:"rating"`
	Servings             uint      `json:"servings"`
	Servings_Size        uint      `json:"servings_size"`
	Prep_Time            uint      `json:"prep_time"`
	Description          string    `json:"description"`
}
type Recipe_Ingredient struct {
	ID                    uint    `json:"id"`
	Ingredient_Mapping_Id uint    `json:"ingredient_mapping_id"`
	Food_Id               uint    `json:"food_id"`
	Amount                float32 `json:"amount"`
	Amount_Unit           string  `json:"amount_unit"`
	Amount_Unit_Desc      string  `json:"amount_unit_desc"`
	Serving_Size          float32 `json:"serving_size"`
}
type Recipe_Review struct {
	ID                 uint      `json:"id"`
	Review_Description string    `json:"review_description"`
	Rating             uint      `json:"rating"`
	Owner_Id           uuid.UUID `json:"owner_id"`
	Food_Id            uint      `json:"food_id"`
	Date_Created       time.Time `json:"date_created"`
}
type Recipe_Likes struct {
	ID           uint      `json:"id"`
	Owner_Id     uint      `json:"owner_id"`
	Date_Created time.Time `json:"date_created"`
}
type Recipe_Instruction struct {
	ID                      uint   `json:"id"`
	Food_Id                 uint   `json:"food_id"`
	Instruction_Description string `json:"instruction_description"`
	Step_Num                uint   `json:"step_num"`
}
type Recipe_Type struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
type Recipe_Image struct {
	ID               uint    `json:"id"`
	Recipe_Id        uint    `json:"recipe_id"`
	Name_File        string  `json:"name_file"`
	Name_URL         string  `json:"name_url"`
	Amount           float32 `json:"amount"`
	Amount_Unit      string  `json:"amount_unit"`
	Amount_Unit_Desc string  `json:"amount_unit_desc"`
}
