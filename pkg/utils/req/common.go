package req

type OTPVerify struct {
	OTP string `json:"otp" binding:"required"`
}
