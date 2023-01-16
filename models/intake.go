package models

import (
	"time"

	"github.com/google/uuid"
)

type Intake struct {
	ID           uint      `json:"id"`
	Account_Id   uuid.UUID `json:"account_id"`
	Food_Id      uint      `json:"food_id"`
	Recipe_Id    uint      `json:"recipe_id"`
	Date_Created time.Time `json:"date_created"`
	// time.Time SHOULD BE IN ISO STRING
}
