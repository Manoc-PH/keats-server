package schemas

// *REQUESTS
type Req_Post_Recipe_Like struct {
	Recipe_Id uint `json:"recipe_id" validate:"required"`
}
