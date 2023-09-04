package schemas

type Req_Post_Thumbnail_Confirm struct {
	ID                   uint   `json:"id" validate:"required"`
	Thumbnail_Image_Link string `json:"thumbnail_image_link" validate:"required"`
}
