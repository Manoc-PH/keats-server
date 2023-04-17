package models

import (
	"time"

	"github.com/google/uuid"
)

type Food struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Name_Ph              string    `json:"name_ph"`
	Name_Brand           string    `json:"name_brand"`
	Date_Created         time.Time `json:"date_created"`
	Thumbnail_Image_Link string    `json:"thumbnail_image_link"`
	Food_Desc            string    `json:"food_desc"`
	Food_Nutrient_Id     uint      `json:"food_nutrient_id"`
	Food_Brand_Type_Id   uint      `json:"food_brand_type_id"`
	Food_Category_Id     uint      `json:"food_category_id"`
	Food_Brand_Id        uuid.UUID `json:"food_brand_id"`
	Removed              bool      `json:"removed"`
	// time.Time SHOULD BE IN ISO STRING
}

type Food_Nutrient struct {
	ID               uint    `json:"id"`
	Food_Id          uint    `json:"food_id"`
	Amount           float32 `json:"amount"`
	Amount_Unit      string  `json:"amount_unit"`
	Amount_Unit_Desc string  `json:"amount_unit_desc"`
	Serving_Size     float32 `json:"serving_size"`
	Calories         float32 `json:"calories"`
	Protein          float32 `json:"protein"`
	Carbs            float32 `json:"carbs"`
	Fats             float32 `json:"fats"`
	//
	Trans_Fat     float32 `json:"trans_fat"`
	Saturated_Fat float32 `json:"saturated_fat"`
	Sugars        float32 `json:"sugars"`
	Sodium        float32 `json:"sodium"`
}
type Food_Image struct {
	ID               uint    `json:"id"`
	Food_Id          uint    `json:"food_id"`
	Name_File        string  `json:"name_file"`
	Amount           float32 `json:"amount"`
	Amount_Unit      string  `json:"amount_unit"`
	Amount_Unit_Desc string  `json:"amount_unit_desc"`
}
type Food_Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type Food_Rewards struct {
	ID      uint `json:"id"`
	Food_Id uint `json:"food_id"`
	Coins   int  `json:"coins"`
	XP      int  `json:"xp"`
}
type Food_Brand struct {
	ID                   uuid.UUID `json:"id"`
	Name                 string    `json:"name"`
	Brand_Desc           string    `json:"brand_desc"`
	Thumbnail_Image_Link string    `json:"thumbnail_image_link"`
	Cover_Image_Link     string    `json:"cover_image_link"`
	Profile_Image_Link   string    `json:"profile_image_link"`
	Food_Brand_Type_Id   uint      `json:"food_brand_type_id"`
}
type Food_Brand_Type struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Brand_Type_Desc string    `json:"brand_type_desc"`
}
