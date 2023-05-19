package models

import (
	"time"
)

type Ingredient struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Name_Ph              string    `json:"name_ph"`
	Name_Owner           string    `json:"name_owner"`
	Date_Created         time.Time `json:"date_created"`
	Thumbnail_Image_Link string    `json:"thumbnail_image_link"`
	Ingredient_Desc      string    `json:"ingredient_desc"`
	Category_Id          uint      `json:"category_id"`
	// time.Time SHOULD BE IN ISO STRING
}
type Ingredient_Variant struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Name_Ph string `json:"name_ph"`
}
type Ingredient_Subvariant struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Name_Ph string `json:"name_ph"`
}
type Ingredient_Mapping struct {
	ID                       uint `json:"id"`
	Ingredient_Id            uint `json:"ingredient_id"`
	Ingredient_Variant_Id    uint `json:"ingredient_variant_id"`
	Ingredient_Subvariant_Id uint `json:"ingredient_subvariant_id"`
	Nutrient_Id              uint `json:"nutrient_id"`
}
type Ingredient_Image struct {
	ID                    uint    `json:"id"`
	Ingredient_Mapping_Id uint    `json:"ingredient_mapping_id"`
	Name_File             string  `json:"name_file"`
	Amount                float32 `json:"amount"`
	Amount_Unit           string  `json:"amount_unit"`
	Amount_Unit_Desc      string  `json:"amount_unit_desc"`
}
