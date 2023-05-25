package interfaces

import "context"

type OTPUseCase interface {
	TwilioSendOtp(context.Context, string) (string, error)
}
