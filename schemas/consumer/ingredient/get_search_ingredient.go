package schemas

// *REQUESTS
type Req_Get_Search_Ingredient struct {
	Search_Term string `json:"search_term" validate:"required"`
}
type Res_Get_Search_Ingredient struct {
	ID                    uint   `json:"id"`
	Name                  string `json:"name"`
	Name_Ph               string `json:"name_ph"`
	Name_Owner            string `json:"name_owner"`
	Thumbnail_Image_Link  string `json:"thumbnail_image_link"`
	Ingredient_Mapping_ID uint   `json:"ingredient_mapping_id"`
}
