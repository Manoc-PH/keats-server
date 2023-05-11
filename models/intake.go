package models

import (
	"time"

	"github.com/google/uuid"
)

type Intake struct {
	ID                    uint      `json:"id"`
	Account_Id            uuid.UUID `json:"account_id"`
	Date_Created          time.Time `json:"date_created"`
	Amount                float32   `json:"amount"`
	Amount_Unit           string    `json:"amount_unit"`
	Amount_Unit_Desc      string    `json:"amount_unit_desc"`
	Serving_Size          float32   `json:"serving_size"`
	Food_Id               uint      `json:"food_id"`
	Ingredient_Mapping_Id uint      `json:"ingredient_mapping_id"`
	// time.Time SHOULD BE IN ISO STRING
}
