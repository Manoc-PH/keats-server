package schemas

// *REQUESTS
type Req_Delete_Recipe struct {
	ID uint `json:"id" validate:"required"`
}
