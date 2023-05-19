package model

type NewAdminInfo struct {
	UserName string `json:"user_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminLoginInfo struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AdminDataOutput struct {
	ID           uint
	UserName     string
	Email        string
	Phone        string
	IsSuperAdmin bool
}
