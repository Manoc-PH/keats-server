package schemas

import "github.com/google/uuid"

type Req_Patch_Main_Image struct {
	Recipe_Id           uuid.UUID `json:"recipe_id" validate:"required"`
	Thumbnail_Name      string    `json:"thumbnail_name"`
	Thumbnail_URL       string    `json:"thumbnail_url"`
	Thumbnail_URL_Local string    `json:"thumbnail_url_local"`
	Image_Name          string    `json:"image_name"`
	Image_URL           string    `json:"image_url"`
	Image_URL_Local     string    `json:"image_url_local"`
}
type Res_Patch_Main_Image struct {
	Recipe_Id           uuid.UUID `json:"recipe_id" validate:"required"`
	Thumbnail_Name      string    `json:"thumbnail_name"`
	Thumbnail_URL       string    `json:"thumbnail_url"`
	Thumbnail_URL_Local string    `json:"thumbnail_url_local"`
	Image_Name          string    `json:"image_name"`
	Image_URL           string    `json:"image_url"`
	Image_URL_Local     string    `json:"image_url_local"`
	Image_Signature     string    `json:"image_signature"`
	Thumbnail_Signature string    `json:"thumbnail_signature"`
	Timestamp           string    `json:"timestamp"`
	Upload_URL          string    `json:"upload_url"`
	API_key             string    `json:"api_key"`
}
