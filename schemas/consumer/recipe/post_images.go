package schemas

type Req_Post_Images struct {
	Recipe_Image []Recipe_Image_Schema `json:"recipe_images" validate:"required,max=7,dive"`
}
type Res_Post_Images struct {
	Recipe_Image []Recipe_Image_Schema `json:"recipe_images"`
}

type Recipe_Image_Schema struct {
	ID               uint    `json:"id"`
	Recipe_Id        uint    `json:"recipe_id" validate:"required"`
	Name_File        string  `json:"name_file" validate:"required"`
	Name_URL         string  `json:"name_url" validate:"required"`
	Amount           float32 `json:"amount"`
	Amount_Unit      string  `json:"amount_unit"`
	Amount_Unit_Desc string  `json:"amount_unit_desc"`
}
