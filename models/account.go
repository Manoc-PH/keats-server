package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Password     []byte    `json:"-"` // putting a minus means not returning field
	Date_Updated time.Time `json:"date_updated"`
	Date_Created time.Time `json:"created"`
	// profile
	Profile_Image_Link string `json:"profile_image_link"`
	Profile_Title      string `json:"profile_title"`
	// info
	Measure_Unit_Id uint      `json:"measure_unit_id"`
	Weight          uint      `json:"weight"`
	Height          uint      `json:"height"`
	Birthday        time.Time `json:"birthday"`
	Sex             string    `json:"sex"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id"`
}
