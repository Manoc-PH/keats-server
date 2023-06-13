package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	Password        []byte    `json:"-"` // putting a minus means not returning field
	Account_Type_Id uuid.UUID `json:"account_type_id"`
}

type Account_Type struct {
	Id                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Account_Type_Desc string    `json:"account_type_desc"`
}

type Consumer_Vitals struct {
	ID              uuid.UUID `json:"id"`
	Account_Id      uuid.UUID `json:"account_id"`
	Weight          int       `json:"weight"`
	Height          int       `json:"height"`
	Birthday        time.Time `json:"birthday"`
	Sex             string    `json:"sex"`
	Activity_Lvl_Id uuid.UUID `json:"activity_lvl_id"`
	Diet_Plan_Id    uuid.UUID `json:"diet_plan_id"`
	// time.Time SHOULD BE IN ISO STRING
}

type Consumer_Weight_Changes struct {
	ID           uuid.UUID `json:"id"`
	Account_Id   uuid.UUID `json:"account_id"`
	Weight       int       `json:"weight"`
	Date_Created time.Time `json:"date_created"`
}

type Consumer_Profile struct {
	ID                 uuid.UUID `json:"id"`
	Account_Id         uuid.UUID `json:"account_id"`
	Account_Image_Link string    `json:"account_image_link"`
	Name_First         string    `json:"name_first"`
	Name_Last          string    `json:"name_last"`
	Phone_Number       string    `json:"phone_number"`
	Date_Updated       time.Time `json:"date_updated"`
	Date_Created       time.Time `json:"date_created"`
	Account_Vitals_Id  uuid.UUID `json:"account_vitals_id"`
	Measure_Unit_Id    uuid.UUID `json:"measure_unit_id"`
	// time.Time SHOULD BE IN ISO STRING
}
type Business_Profile struct {
	ID                 uuid.UUID `json:"id"`
	Account_Id         uuid.UUID `json:"account_id"`
	Account_Image_Link string    `json:"account_image_link"`
	Name_Business      string    `json:"name_business"`
	Phone_Number       string    `json:"phone_number"`
	Date_Updated       time.Time `json:"date_updated"`
	Date_Created       time.Time `json:"date_created"`
	// time.Time SHOULD BE IN ISO STRING
}
