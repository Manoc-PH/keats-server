package schemas

// *REQUESTS
type Req_Delete_Recipe_Review struct {
	Recipe_ID uint `json:"recipe_id" validate:"required"`
}
