package schemas

import (
	"server/models"

	"github.com/google/uuid"
)

// *REQUESTS
type Req_Get_Food_Details struct {
	Food_ID uint `json:"food_id" validate:"required"`
}

type Res_Get_Food_Details struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Name_Ph              string    `json:"name_ph"`
	Name_Brand           string    `json:"name_brand"`
	Food_Desc            string    `json:"food_desc"`
	Barcode              string    `json:"barcode"`
	Thumbnail_Image_Link string    `json:"thumbnail_image_link"`
	Food_Nutrient_Id     uint      `json:"food_nutrient_id"`
	Food_Brand_Type_Id   uint      `json:"food_brand_type_id"`
	Food_Category_Id     uint      `json:"food_category_id"`
	Food_Brand_Id        uuid.UUID `json:"food_brand_id"`
	//
	Food_Nutrients  models.Food_Nutrient   `json:"food_nutrients"`
	Food_Brand_Type models.Food_Brand_Type `json:"food_brand_type"`
	Food_Brand      models.Food_Brand      `json:"food_brand"`
	Food_Category   models.Food_Category   `json:"food_category"`
	Food_Images     []models.Food_Image    `json:"food_images"`
}
