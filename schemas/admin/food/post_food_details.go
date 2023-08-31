package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Food struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name" validate:"required"`
	Name_Ph              string    `json:"name_ph"`
	Name_Owner           string    `json:"name_owner"`
	Date_Created         time.Time `json:"date_created"`
	Barcode              string    `json:"barcode" validate:"required"`
	Thumbnail_Image_Link string    `json:"thumbnail_image_link"`
	Food_Desc            string    `json:"food_desc"`
	Owner_Id             uuid.UUID `json:"owner_id"`
	Category_Id          uint      `json:"category_id"`
	Food_Type_Id         uint      `json:"food_type_id"`
	Nutrient_Id          uint      `json:"nutrient_id"`
	// time.Time SHOULD BE IN ISO STRING
}

type Nutrient struct {
	ID               uint    `json:"id"`
	Amount           float32 `json:"amount" validate:"required"`
	Amount_Unit      string  `json:"amount_unit" validate:"required"`
	Amount_Unit_Desc string  `json:"amount_unit_desc" validate:"required"`
	Serving_Size     float32 `json:"serving_size" validate:"required"`
	Serving_Total    float32 `json:"serving_total" validate:"required"`
	Calories         float32 `json:"calories" validate:"required"`
	Protein          float32 `json:"protein" validate:"required"`
	Carbs            float32 `json:"carbs" validate:"required"`
	Fats             float32 `json:"fats" validate:"required"`
	//
	Trans_Fat     float32 `json:"trans_fat"`
	Saturated_Fat float32 `json:"saturated_fat"`
	Sugars        float32 `json:"sugars"`
	Fiber         float32 `json:"fiber"`
	Sodium        float32 `json:"sodium"`
	Iron          float32 `json:"iron"`
	Calcium       float32 `json:"calcium"`
}
type Req_Post_Food_Details struct {
	Food     Food     `json:"food" validate:"required"`
	Nutrient Nutrient `json:"nutrient" validate:"required"`
}
