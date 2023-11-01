package schemas

import (
	"server/models"

	"github.com/google/uuid"
)

type Req_Get_Recipe_Instructions struct {
	Recipe_Id uuid.UUID `json:"recipe_id" validate:"required"`
}
type Res_Get_Recipe_Instructions struct {
	Instructions []models.Recipe_Instruction `json:"instructions"`
}
