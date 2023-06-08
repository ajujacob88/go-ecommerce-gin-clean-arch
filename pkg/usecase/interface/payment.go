package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
)

type PaymentUseCase interface {
	GetPaymentMethodInfoByID(ctx context.Context, paymentMethodID int) (domain.PaymentMethodInfo, error)

	RazorPayCheckout(ctx context.Context, userID, orderID int) (domain.Order, string, error)

	UpdatePaymentDetails(ctx context.Context, paymentVerifier request.PaymentVerification) error
}
