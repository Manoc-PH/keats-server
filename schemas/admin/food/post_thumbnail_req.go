package schemas

import (
	"time"
)

type Req_Post_Thumbnail_Req struct {
	Timestamp time.Time `json:"timestamp" validate:"required"`
}
type Res_Post_Thumbnail_Req struct {
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp"`
}
