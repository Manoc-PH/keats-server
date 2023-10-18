package schemas

import "server/models"

type Req_Get_Recipe_Details struct {
	Recipe_Id uint `json:"recipe_id" validate:"required"`
}
type Res_Get_Recipe_Details struct {
	Recipe             models.Recipe              `json:"recipe"`
	Recipe_Ingredients []models.Recipe_Ingredient `json:"recipe_ingredients"`
	Recipe_Images      []models.Recipe_Image      `json:"recipe_images"`
	Nutrients          models.Nutrient            `json:"nutrients"`
}
