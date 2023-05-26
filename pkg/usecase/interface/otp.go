package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
)

type OTPUseCase interface {
	TwilioSendOtp(ctx context.Context, phoneNumber string) (string, error)
	TwilioVerifyOTP(ctx context.Context, otpverify model.OTPVerify) (domain.OTPSession, error)
}
