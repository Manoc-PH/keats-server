package constants

type Account_Types_Struct struct {
	Admin    string
	Business string
	Consumer string
}

var Account_Types = Account_Types_Struct{
	Admin:    "admin",
	Business: "business",
	Consumer: "consumer",
}
