package models

import (
	"github.com/google/uuid"
)

type Account_Profile struct {
	ID                 uuid.UUID `json:"id"`
	Account_Id         uuid.UUID `json:"account_id"`
	Profile_Image_Link string    `json:"profile_image_link"`
	Profile_Title      string    `json:"profile_title"`
}
