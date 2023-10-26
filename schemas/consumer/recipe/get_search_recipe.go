package schemas

// *REQUESTS
type Req_Get_Search_Recipe struct {
	Search_Term string `json:"search_term" validate:"required"`
}

type Search_Recipe_Result struct {
	ID                   uint    `json:"id"`
	Name                 string  `json:"name"`
	Name_Ph              string  `json:"name_ph"`
	Name_Owner           string  `json:"name_owner"`
	Thumbnail_Image_Link string  `json:"thumbnail_image_link"`
	Main_Image_Link      string  `json:"main_image_link"`
	Rating               float32 `json:"rating"`
	Rating_Count         uint    `json:"rating_count"`
}

type Meili_Res struct {
	Estimated_Total_Hits int                    `json:"estimatedTotalHits"`
	Limit                int                    `json:"limit"`
	ProcessingTimeMs     int                    `json:"processingTimeMs"`
	Query                string                 `json:"query"`
	Hits                 []Search_Recipe_Result `json:"hits"`
}
