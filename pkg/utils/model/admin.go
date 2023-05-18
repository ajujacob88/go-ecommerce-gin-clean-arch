package model

type NewAdminInfo struct {
	UserName string `json:"user_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"enail" validate:"required"`
	Password string `json:"password" validate:"required"`
}
