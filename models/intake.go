package models

import (
	"time"

	"github.com/google/uuid"
)

type Intake struct {
	ID           int       `json:"id"`
	Account_Id   uuid.UUID `json:"account_id"`
	Food_Id      int       `json:"food_id"`
	Recipe_Id    int       `json:"recipe_id"`
	Date_Created time.Time `json:"date_created"`
}
