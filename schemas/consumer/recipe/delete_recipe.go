package schemas

import "github.com/google/uuid"

// *REQUESTS
type Req_Delete_Recipe struct {
	ID uuid.UUID `json:"id" validate:"required"`
}
