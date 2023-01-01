package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password []byte    `json:"-"` // putting a minus means not returning field
	Updated  time.Time `json:"updated"`
	Created  time.Time `json:"created"`
	// profile
	Profile_Image_Link string `json:"profile_image_link"`
	Profile_Title      string `json:"profile_title"`
	// info
	Weight          uint   `json:"weight"`
	Height          uint   `json:"height"`
	Age             uint   `json:"age"`
	Sex             string `json:"sex"`
	Activity_Lvl_Id uint   `json:"activity_lvl_id"`
	Diet_Plan_Id    uint   `json:"diet_plan_id"`
}
