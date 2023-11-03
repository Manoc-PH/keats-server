package schemas

import (
	"server/models"
)

type Res_Get_Recipe_Discovery struct {
	Recipes []models.Recipe `json:"recipes"`
}
