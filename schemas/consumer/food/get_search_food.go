package schemas

import "github.com/google/uuid"

// *REQUESTS
type Req_Get_Search_Food struct {
	Search_Term string `json:"search_term" validate:"required"`
}
type Res_Get_Search_Food struct {
	ID                   uuid.UUID `json:"id"`
	Name                 string    `json:"name"`
	Name_Ph              string    `json:"name_ph"`
	Name_Brand           string    `json:"name_brand"`
	Thumbnail_Image_Link string    `json:"thumbnail_image_link"`
	Food_Nutrient_Id     uuid.UUID `json:"food_nutrient_id"`
	Ranking              float32   `json:"ranking"`
}
