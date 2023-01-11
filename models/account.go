package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	Password           []byte    `json:"-"` // putting a minus means not returning field
	Date_Updated       time.Time `json:"date_updated"`
	Date_Created       time.Time `json:"created"`
	Account_Vitals_Id  uuid.UUID `json:"account_vitals_id"`
	Account_profile_Id uuid.UUID `json:"account_profile_id"`
	Measure_Unit_Id    uuid.UUID `json:"measure_unit_id"`
}
