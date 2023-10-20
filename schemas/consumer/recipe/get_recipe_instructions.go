package schemas

import "server/models"

type Req_Get_Recipe_Instructions struct {
	Recipe_Id uint `json:"recipe_id" validate:"required"`
}
type Res_Get_Recipe_Instructions struct {
	Instructions []models.Recipe_Instruction `json:"instructions"`
}
