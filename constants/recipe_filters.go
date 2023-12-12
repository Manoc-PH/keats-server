package constants

type Recipe_Nutrition_Category_Struct struct {
	H_Protein string
	H_Carbs   string
	H_Fats    string
	L_Cal     string
	L_Carbs   string
	L_Fats    string
}
type Recipe_Nutrition_Category_Order_Struct struct {
	DESC string
	ASC  string
}

var Recipe_Nutrition_Category_Order = Recipe_Nutrition_Category_Order_Struct{
	DESC: "DESC",
	ASC:  "ASC",
}
var Recipe_Nutrition_Categories = Recipe_Nutrition_Category_Struct{
	H_Protein: "h_protein",
	H_Carbs:   "h_carbs",
	H_Fats:    "h_fats",
	L_Cal:     "l_cal",
	L_Carbs:   "l_carbs",
	L_Fats:    "l_fats",
}

func Get_Nutrition_Categories(v string) string {
	var Recipe_Nutrition_Categories_SQL = map[string]string{
		"h_protein": "protein",
		"h_carbs":   "carbs",
		"h_fats":    "fats",
		"l_cal":     "calories",
		"l_carbs":   "carbs",
		"l_fats":    "fats",
	}

	value, exists := Recipe_Nutrition_Categories_SQL[v]
	if exists {
		return value
	} else {
		return ""
	}
}
func Get_Nutrition_Categories_Order(v string) string {
	var Recipe_Nutrition_Categories_Order = map[string]string{
		"h_protein": "DESC",
		"h_carbs":   "DESC",
		"h_fats":    "DESC",
		"l_cal":     "ASC",
		"l_carbs":   "ASC",
		"l_fats":    "ASC",
	}
	value, exists := Recipe_Nutrition_Categories_Order[v]
	if exists {
		return value
	} else {
		return ""
	}
}
