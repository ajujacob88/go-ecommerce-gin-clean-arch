package interfaces

import "context"

type OTPRepository interface {
	SaveOTP(ctx context.Context, resp string, phoneNumber string) error
}
