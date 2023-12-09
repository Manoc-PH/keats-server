package schemas

import "github.com/google/uuid"

type Req_Post_Image struct {
	ID               uuid.UUID `json:"id"`
	Recipe_Id        uuid.UUID `json:"recipe_id" validate:"required"`
	Name_File        string    `json:"name_file"`
	Name_URL         string    `json:"name_url"`
	Name_URL_Local   string    `json:"name_url_local"`
	Amount           float32   `json:"amount"`
	Amount_Unit      string    `json:"amount_unit"`
	Amount_Unit_Desc string    `json:"amount_unit_desc"`
}
type Res_Post_Image struct {
	ID               uuid.UUID `json:"id"`
	Recipe_Id        uuid.UUID `json:"recipe_id" validate:"required"`
	Name_File        string    `json:"name_file"`
	Name_URL         string    `json:"name_url"`
	Name_URL_Local   string    `json:"name_url_local"`
	Amount           float32   `json:"amount"`
	Amount_Unit      string    `json:"amount_unit"`
	Amount_Unit_Desc string    `json:"amount_unit_desc"`
	Signature        string    `json:"signature"`
	Timestamp        string    `json:"timestamp"`
	Upload_URL       string    `json:"upload_url"`
	API_key          string    `json:"api_key"`
}
