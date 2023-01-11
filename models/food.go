package models

import "time"

type Food struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Name_Ph          string    `json:"name_ph"`
	Name_Brand       string    `json:"name_brand"`
	Food_Nutrient_Id int       `json:"food_nutrient_id"`
	Date_Created     time.Time `json:"date_created"`
}
