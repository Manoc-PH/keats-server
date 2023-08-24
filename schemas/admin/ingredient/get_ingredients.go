package schemas

type Req_Get_Ingredients struct {
	Index uint `json:"index" validate:"min=0"`
	Size  uint `json:"size" validate:"required,min=10,max=30"`
}
