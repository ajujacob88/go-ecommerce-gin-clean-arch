package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

type PaymentRepository interface {
	GetPaymentMethodInfoByID(ctx context.Context, paymentMethodID int) (domain.PaymentMethodInfo, error)
}
