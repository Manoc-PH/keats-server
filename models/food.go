package models

import (
	"time"

	"github.com/google/uuid"
)

type Food struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Name_Ph              string    `json:"name_ph"`
	Name_Owner           string    `json:"name_owner"`
	Date_Created         time.Time `json:"date_created"`
	Barcode              string    `json:"barcode"`
	Thumbnail_Image_Link string    `json:"thumbnail_image_link"`
	Food_Desc            string    `json:"food_desc"`
	Category_Id          uint      `json:"category_id"`
	Food_Type_Id         uint      `json:"food_type_id"`
	Owner_Id             uuid.UUID `json:"owner_id"`
	// time.Time SHOULD BE IN ISO STRING
}

type Food_Image struct {
	ID               uint    `json:"id"`
	Food_Id          uint    `json:"food_id"`
	Name_File        string  `json:"name_file"`
	Amount           float32 `json:"amount"`
	Amount_Unit      string  `json:"amount_unit"`
	Amount_Unit_Desc string  `json:"amount_unit_desc"`
}

type Food_Ingredient struct {
	ID                    uint    `json:"id"`
	Ingredient_Mapping_Id uint    `json:"ingredient_mapping_id"`
	Amount                float32 `json:"amount"`
	Amount_Unit           string  `json:"amount_unit"`
	Amount_Unit_Desc      string  `json:"amount_unit_desc"`
	Serving_Size          float32 `json:"serving_size"`
}

// This table will be used in the case where the food does not have any ingredients.
// Only one food_nutrient exists per food
type Food_Nutrient struct {
	ID          uint `json:"id"`
	Nutrient_Id uint `json:"nutrient_id"`
	Food_Id     uint `json:"food_id"`
}

type Edible_Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type Food_Type struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
