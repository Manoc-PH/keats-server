package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Req_Get_Recipe_Review struct {
	Recipe_Id uuid.UUID `json:"recipe_id" validate:"required"`
}
type Res_Get_Recipe_Review struct {
	ID           uuid.UUID `json:"id"`
	Description  string    `json:"description"`
	Rating       float32   `json:"rating"`
	Owner_Id     uuid.UUID `json:"owner_id"`
	Recipe_Id    uuid.UUID `json:"recipe_id"`
	Date_Created time.Time `json:"date_created"`
	Name_Owner   string    `json:"name_owner"`
}
