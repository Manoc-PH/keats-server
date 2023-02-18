package schemas

import (
	"server/models"

	"github.com/google/uuid"
)

// *REQUESTS
type Req_Get_Food_Details struct {
	Food_ID uint `json:"food_id" validate:"required"`
}

type Food_Details struct {
	ID                 uint      `json:"id"`
	Name               string    `json:"name"`
	Name_Ph            string    `json:"name_ph"`
	Name_Brand         string    `json:"name_brand"`
	Food_Desc          string    `json:"food_desc"`
	Barcode            string    `json:"barcode"`
	Food_Nutrient_Id   uint      `json:"food_nutrient_id"`
	Food_Brand_Type_Id uint      `json:"food_brand_type_id"`
	Food_Category_Id   uint      `json:"food_category_id"`
	Food_Brand_Id      uuid.UUID `json:"food_brand_id"`
	//
	Amount           float32 `json:"amount"`
	Amount_Unit      string  `json:"amount_unit"`
	Amount_Unit_Desc string  `json:"amount_unit_desc"`
	Serving_Size     float32 `json:"serving_size"`
	Calories         float32 `json:"calories"`
	Protein          float32 `json:"protein"`
	Carbs            float32 `json:"carbs"`
	Fats             float32 `json:"fats"`
	Trans_Fat        float32 `json:"trans_fat"`
	Saturated_Fat    float32 `json:"saturated_fat"`
	Sugars           float32 `json:"sugars"`
	Sodium           float32 `json:"sodium"`
}
type Res_Get_Food_Details struct {
	Food_Details Food_Details        `json:"food_details"`
	Food_Images  []models.Food_Image `json:"food_images"`
}
