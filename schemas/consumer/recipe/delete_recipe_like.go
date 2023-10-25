package schemas

// *REQUESTS
type Req_Delete_Recipe_Like struct {
	Recipe_ID uint `json:"recipe_id" validate:"required"`
}
