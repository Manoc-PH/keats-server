package schemas

import "github.com/google/uuid"

type Req_Login struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=10,max=32"` //,missingRequiredCharacters  // add this for password validation
}
type Res_Login struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	Account_Type_Id uuid.UUID `json:"account_type_id"`
	Token           string    `json:"token"`
}
