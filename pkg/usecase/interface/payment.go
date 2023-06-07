package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

type PaymentUseCase interface {
	GetPaymentMethodInfoByID(ctx context.Context, paymentMethodID int) (domain.PaymentMethodInfo, error)

	RazorPayCheckout(ctx context.Context, userID, orderID int) (domain.Order, string, error)
}
