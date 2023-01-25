package schemas

// *REQUESTS
type Req_Get_Search_Food struct {
	Search_Term string `json:"search_term" validate:"required"`
}
