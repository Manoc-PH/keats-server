package models

type Game_Item struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	Price                uint   `json:"price"`
	Commentary           string `json:"commentary"`
	Game_Item_Desc       string `json:"game_item_desc"`
	Game_Item_Image_Link string `json:"game_item_image_link"`
}
