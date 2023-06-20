package schemas

type Req_Get_Name_Availability struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
}
