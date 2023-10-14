package models

import (
	"time"

	"github.com/google/uuid"
)

type Food_Ingredient struct {
	ID                    uint    `json:"id"`
	Ingredient_Mapping_Id uint    `json:"ingredient_mapping_id"`
	Food_Id               uint    `json:"food_id"`
	Amount                float32 `json:"amount"`
	Amount_Unit           string  `json:"amount_unit"`
	Amount_Unit_Desc      string  `json:"amount_unit_desc"`
	Serving_Size          float32 `json:"serving_size"`
}

type Food_Recipe_Details struct {
	ID               uint   `json:"id"`
	Food_Id          uint   `json:"food_id"`
	Main_Image_Link  string `json:"main_image_link"`
	Likes            uint   `json:"likes"`
	Servings         uint   `json:"servings"`
	Servings_Size    uint   `json:"servings_size"`
	Prep_Time        uint   `json:"prep_time"`
	Food_Description string `json:"food_description"`
}
type Food_Review struct {
	ID                 uint      `json:"id"`
	Review_Description string    `json:"review_description"`
	Rating             uint      `json:"rating"`
	Owner_Id           uuid.UUID `json:"owner_id"`
	Food_Id            uint      `json:"food_id"`
	Date_Created       time.Time `json:"date_created"`
}
type Food_Likes struct {
	ID           uint      `json:"id"`
	Owner_Id     uint      `json:"owner_id"`
	Date_Created time.Time `json:"date_created"`
}
type Food_Instruction struct {
	ID                      uint   `json:"id"`
	Food_Id                 uint   `json:"food_id"`
	Instruction_Description string `json:"instruction_description"`
	Step_Num                uint   `json:"step_num"`
}
type Food_Type struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
