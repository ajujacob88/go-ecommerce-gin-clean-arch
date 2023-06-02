package response

//model structs used to output data

type AdminDataOutput struct {
	ID           uint
	UserName     string
	Email        string
	Phone        string
	IsSuperAdmin bool
}
