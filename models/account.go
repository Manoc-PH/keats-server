package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	Name_First         string    `json:"name_first"`
	Name_Last          string    `json:"name_last"`
	Phone_Number       string    `json:"phone_number"`
	Password           []byte    `json:"-"` // putting a minus means not returning field
	Date_Updated       time.Time `json:"date_updated"`
	Date_Created       time.Time `json:"created"`
	Account_Type_Id    uuid.UUID `json:"account_type_id"`
	Account_Vitals_Id  uuid.UUID `json:"account_vitals_id"`
	Account_profile_Id uuid.UUID `json:"account_profile_id"`
	Measure_Unit_Id    uuid.UUID `json:"measure_unit_id"`
	// time.Time SHOULD BE IN ISO STRING
}

type Account_Type struct {
	Id                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Account_Type_Desc string    `json:"account_type_desc"`
}

type Account_Vitals struct {
	ID         uuid.UUID `json:"id"`
	Account_Id uuid.UUID `json:"account_id"`
	// info
	Weight          int       `json:"weight"`
	Height          int       `json:"height"`
	Birthday        time.Time `json:"birthday"`
	Sex             string    `json:"sex"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id"`
	// time.Time SHOULD BE IN ISO STRING
}

type Account_Profile struct {
	ID                 uuid.UUID `json:"id"`
	Account_Id         uuid.UUID `json:"account_id"`
	Profile_Image_Link string    `json:"profile_image_link"`
	Profile_Title      string    `json:"profile_title"`
}
type Account_Weight_Changes struct {
	ID           uuid.UUID `json:"id"`
	Account_Id   uuid.UUID `json:"account_id"`
	Weight       int       `json:"weight"`
	Date_Created time.Time `json:"created"`
}

type Account_Items struct {
	ID           uint      `json:"id"`
	Account_Id   uuid.UUID `json:"account_id"`
	Game_Item_Id uint      `json:"game_item_id"`
}

type Account_Game_Stat struct {
	ID         uint      `json:"id"`
	Account_Id uuid.UUID `json:"account_id"`
	Coins      uint      `json:"coins"`
	XP         uint      `json:"xp"`
}
