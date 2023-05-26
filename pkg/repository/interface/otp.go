package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
)

type OTPRepository interface {
	SaveOTP(ctx context.Context, resp string, phoneNumber string) error
	RetrieveOtpSession(ctx context.Context, otpverify model.OTPVerify) (domain.OTPSession, error)
}
