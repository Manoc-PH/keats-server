package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Req_Get_Recipe_Reviews struct {
	Recipe_Id uint `json:"recipe_id" validate:"required"`
	Page      uint `json:"page" validate:"min=0"`
	Size      uint `json:"size" validate:"required,min=1,max=40"`
}
type Res_Get_Recipe_Reviews struct {
	Reviews []Recipe_Reviews_Schema `json:"reviews"`
}
type Recipe_Reviews_Schema struct {
	ID           uint      `json:"id"`
	Description  string    `json:"description"`
	Rating       float32   `json:"rating"`
	Owner_Id     uuid.UUID `json:"owner_id"`
	Recipe_Id    uint      `json:"recipe_id"`
	Date_Created time.Time `json:"date_created"`
	Name_Owner   string    `json:"name_owner"`
}
