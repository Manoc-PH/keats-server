package schemas

import "github.com/google/uuid"

// *REQUESTS
type Req_Post_Recipe_Like struct {
	Recipe_Id uuid.UUID `json:"recipe_id" validate:"required"`
}
