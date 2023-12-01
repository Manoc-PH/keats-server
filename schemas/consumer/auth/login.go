package schemas

type Req_Login struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=10,max=32"` //,missingRequiredCharacters  // add this for password validation
}
type Res_Login struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
