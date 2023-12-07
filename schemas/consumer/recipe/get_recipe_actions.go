package schemas

import (
	"github.com/google/uuid"
)

type Req_Get_Actions struct {
	Recipe_Id uuid.UUID `json:"recipe_id" validate:"required"`
}
type Res_Get_Actions struct {
	Liked    bool `json:"liked"`
	Reviewed bool `json:"reviewed"`
}
