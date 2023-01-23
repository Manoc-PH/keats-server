package schemas

import (
	"server/models"
	"time"

	"github.com/google/uuid"
)

type Intake_Details_Food struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Name_Ph          string `json:"name_ph"`
	Name_Brand       string `json:"name_brand"`
	Food_Nutrient_Id uint   `json:"food_nutrient_id"`
	// Food nutrient
	Amount           float32 `json:"amount"`
	Amount_Unit      string  `json:"amount_unit"`
	Amount_Unit_Desc string  `json:"amount_unit_desc"`
	Serving_Size     float32 `json:"serving_size"`
	Calories         float32 `json:"calories"`
	Protein          float32 `json:"protein"`
	Carbs            float32 `json:"carbs"`
	Fats             float32 `json:"fats"`
	//
	Trans_Fat     float32 `json:"trans_fat"`
	Saturated_Fat float32 `json:"saturated_fat"`
	Sugars        float32 `json:"sugars"`
	Sodium        float32 `json:"sodium"`
}
type Intake_Food struct {
	Details Intake_Details_Food `json:"details"`
	Images  []models.Food_Image `json:"food_image"`
}
type Intake_Recipe struct {
}
type Res_Get_Intake_Details struct {
	ID               uint          `json:"id"`
	Account_Id       uuid.UUID     `json:"account_id"`
	Food_Id          uint          `json:"food_id"`
	Recipe_Id        uint          `json:"recipe_id"`
	Date_Created     time.Time     `json:"date_created"`
	Amount           float32       `json:"amount"`
	Amount_Unit      string        `json:"amount_unit"`
	Amount_Unit_Desc string        `json:"amount_unit_desc"`
	Serving_Size     float32       `json:"serving_size"`
	Food             Intake_Food   `json:"food"`
	Recipe           Intake_Recipe `json:"recipe"`
}
