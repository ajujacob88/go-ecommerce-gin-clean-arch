package request

//model structs input

type OTPVerify struct {
	OTP   string `json:"otp" binding:"required"`
	OtpId string `json:"otpid" binding:"required"`
}
