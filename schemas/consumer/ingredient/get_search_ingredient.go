package schemas

// *REQUESTS
type Req_Get_Search_Ingredient struct {
	Search_Term string `json:"search_term" validate:"required"`
}

type Ingredient_Details struct {
	Ingredient_Mapping_Id      uint   `json:"ingredient_mapping_id"`
	Ingredient_Variant_Name    string `json:"ingredient_variant_name"`
	Ingredient_Subvariant_Name string `json:"ingredient_subvariant_name"`
}
type Search_Ingredient_Result struct {
	ID                   uint               `json:"id"`
	Name                 string             `json:"name"`
	Name_Owner           string             `json:"name_owner"`
	Thumbnail_Image_Link string             `json:"thumbnail_image_link"`
	Ingredient_Details   Ingredient_Details `json:"ingredient_details"`
}

type Meili_Res struct {
	Estimated_Total_Hits int                        `json:"estimatedTotalHits"`
	Limit                int                        `json:"limit"`
	ProcessingTimeMs     int                        `json:"processingTimeMs"`
	Query                string                     `json:"query"`
	Hits                 []Search_Ingredient_Result `json:"hits"`
}
