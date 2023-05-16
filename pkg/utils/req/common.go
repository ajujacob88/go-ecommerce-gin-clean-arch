package req

type OTPVerify struct {
	OTP string `json:"otp" binding:"required"`
}

type UserLoginEmail struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5,max=30"`
}
