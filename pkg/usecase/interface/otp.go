package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
)

type OTPUseCase interface {
	TwilioSendOtp(ctx context.Context, phoneNumber string) (string, error)
	TwilioVerifyOTP(ctx context.Context, otpverify request.OTPVerify) (domain.OTPSession, error)
}
