package schemas

import "github.com/google/uuid"

// *REQUESTS
type Req_Delete_Recipe_Review struct {
	Recipe_ID uuid.UUID `json:"recipe_id" validate:"required"`
}
