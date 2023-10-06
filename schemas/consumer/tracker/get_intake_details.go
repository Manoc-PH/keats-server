package schemas

import (
	"server/models"
	"time"

	"github.com/google/uuid"
)

type Req_Get_Intake_Details struct {
	Intake_ID uint `json:"intake_id" validate:"required"`
}
type Intake_Ingredient struct {
	Details Ingredient_Mapping_Schema `json:"details"`
	Images  []models.Ingredient_Image `json:"images"`
}
type Intake_Food struct {
	Details Food_Mapping_Schema `json:"details"`
	Images  []models.Food_Image `json:"images"`
}
type Res_Get_Intake_Details struct {
	ID                    uint               `json:"id"`
	Account_Id            uuid.UUID          `json:"account_id"`
	Food_Id               uint               `json:"food_id"`
	Ingredient_Mapping_Id uint               `json:"ingredient_mapping_id"`
	Date_Created          time.Time          `json:"date_created"`
	Amount                float32            `json:"amount"`
	Amount_Unit           string             `json:"amount_unit"`
	Amount_Unit_Desc      string             `json:"amount_unit_desc"`
	Serving_Size          float32            `json:"serving_size"`
	Ingredient            *Intake_Ingredient `json:"ingredient"`
	Food                  *Intake_Food       `json:"food"`
}
