package models

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID              int       `json:"id"`
	Owner_Id        uuid.UUID `json:"owner_id"`
	Main_Image_Link string    `json:"main_image_link"`
	Prep_Mins       int       `json:"prep_mins"`
	Servings        int       `json:"servings"`
	Saves           int       `json:"saves"`
	Date_Created    time.Time `json:"date_created"`
	Date_Updated    time.Time `json:"date_updated"`
}
