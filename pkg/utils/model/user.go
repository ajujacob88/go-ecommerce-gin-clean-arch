package model

type BlockUser struct {
	UserID int    `json:"user_id"`
	Reason string `json:"reason"`
}

type NewUserInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"required, email"`
	Phone     string `json:"phone" binding:"required, ,min=10,max=10"`
	Password  string `json:"password" validate:"required"`
}

type UserDataOutput struct {
	ID        uint   `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
