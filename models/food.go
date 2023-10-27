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
	Description          string    `json:"description"`
	Category_Id          uint      `json:"category_id"`
	Owner_Id             uuid.UUID `json:"owner_id"`
	Nutrient_Id          uint      `json:"nutrient_id"`
	// time.Time SHOULD BE IN ISO STRING
}
type Food_Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type Food_Image struct {
	ID               uint    `json:"id"`
	Food_Id          uint    `json:"food_id"`
	Name_File        string  `json:"name_file"`
	Name_URL         string  `json:"name_url"`
	Amount           float32 `json:"amount"`
	Amount_Unit      string  `json:"amount_unit"`
	Amount_Unit_Desc string  `json:"amount_unit_desc"`
}
