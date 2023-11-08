package schemas

import "github.com/google/uuid"

// *REQUESTS
type Req_Get_Search_Recipe struct {
	Search_Term string `json:"search_term" validate:"required"`
}

type Search_Recipe_Result struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Name_Ph        string    `json:"name_ph"`
	Name_Owner     string    `json:"name_owner"`
	Thumbnail_URL  string    `json:"thumbnail_url"`
	Thumbnail_Name string    `json:"thumbnail_name"`
	Image_URL      string    `json:"image_url"`
	Image_Name     string    `json:"image_name"`
	Rating         float32   `json:"rating"`
	Rating_Count   uint      `json:"rating_count"`
}

type Meili_Res struct {
	Estimated_Total_Hits int                    `json:"estimatedTotalHits"`
	Limit                int                    `json:"limit"`
	ProcessingTimeMs     int                    `json:"processingTimeMs"`
	Query                string                 `json:"query"`
	Hits                 []Search_Recipe_Result `json:"hits"`
}
