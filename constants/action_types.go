package constants

type Action_Types_Struct struct {
	Delete string
	Update string
	Insert string
}

var Action_Types = Action_Types_Struct{
	Delete: "delete",
	Update: "update",
	Insert: "insert",
}
