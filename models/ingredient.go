package models

import (
	"time"

	"github.com/google/uuid"
)

type Ingredient struct {
	ID                   uuid.UUID `json:"id"`
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
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Name_Ph string    `json:"name_ph"`
}
type Ingredient_Subvariant struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Name_Ph string    `json:"name_ph"`
}
type Ingredient_Mapping struct {
	ID                       uuid.UUID `json:"id"`
	Ingredient_Id            uuid.UUID `json:"ingredient_id"`
	Ingredient_Variant_Id    uuid.UUID `json:"ingredient_variant_id"`
	Ingredient_Subvariant_Id uuid.UUID `json:"ingredient_subvariant_id"`
	Nutrient_Id              uuid.UUID `json:"nutrient_id"`
}
type Ingredient_Image struct {
	ID                    uuid.UUID `json:"id"`
	Ingredient_Mapping_Id uuid.UUID `json:"ingredient_mapping_id"`
	Name_File             string    `json:"name_file"`
	Name_URL              string    `json:"name_url"`
	Amount                float32   `json:"amount"`
	Amount_Unit           string    `json:"amount_unit"`
	Amount_Unit_Desc      string    `json:"amount_unit_desc"`
}

type Ingredient_Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
