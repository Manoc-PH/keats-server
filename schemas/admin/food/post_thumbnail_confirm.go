package schemas

import "github.com/google/uuid"

type Req_Post_Thumbnail_Confirm struct {
	ID                   uuid.UUID `json:"id" validate:"required"`
	Thumbnail_Image_Link string    `json:"thumbnail_image_link" validate:"required"`
}
