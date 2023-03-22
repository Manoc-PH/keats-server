package schemas

import "github.com/google/uuid"

// TODO Add birthday and remove age
type Req_Update_Account_Profile struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Username string    `json:"username" validate:"required,min=3,max=32"`

	Profile_Image_Link string `json:"profile_image_link" validate:"required"`
	Profile_Title      string `json:"profile_title" validate:"required"`

	Weight          uint   `json:"weight" validate:"required,min=1,max=200"`
	Height          uint   `json:"height" validate:"required,min=1,max=250"`
	Age             uint   `json:"age" validate:"required,min=16,max=100"`
	Sex             string `json:"sex" validate:"required,min=1,max=1"`
	Activity_Lvl_Id uint   `json:"activity_lvl_id" validate:"required,min=1,max=32"`
	Diet_Plan_Id    uint   `json:"diet_plan_id" validate:"required,min=1,max=32"`
}
