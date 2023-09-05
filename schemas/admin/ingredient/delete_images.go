package schemas

type Ingredient_Image_Schema struct {
	ID       uint   `json:"id" validate:"required"`
	Name_URL string `json:"name_url" validate:"required"`
}

type Req_Delete_Images struct {
	Images []Ingredient_Image_Schema `json:"images" validate:"required"`
}
