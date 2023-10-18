package schemas

import (
	"time"

	"github.com/google/uuid"
)

// *REQUESTS
type Req_Post_Recipe_Review struct {
	ID           uint      `json:"id"`
	Description  string    `json:"description"`
	Rating       float32   `json:"rating" validate:"required,min=1,max=5"`
	Owner_Id     uuid.UUID `json:"owner_id"`
	Recipe_Id    uint      `json:"recipe_id" validate:"required"`
	Date_Created time.Time `json:"date_created"`
}
