package schemas

type Req_Sign_Up struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=10,max=32"` //,missingRequiredCharacters // add this for password validation

	Weight          uint   `json:"weight" validate:"required,min=1,max=200"`
	Height          uint   `json:"height" validate:"required,min=1,max=250"`
	Age             uint   `json:"age" validate:"required,min=16,max=100"`
	Sex             string `json:"sex" validate:"required,min=1,max=1"`
	Activity_Lvl_Id uint   `json:"activity_lvl_id" validate:"required,min=1,max=32"`
	Diet_Plan_Id    uint   `json:"diet_plan_id" validate:"required,min=1,max=32"`
}
