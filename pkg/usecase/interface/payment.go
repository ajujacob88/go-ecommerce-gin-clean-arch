package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
)

type PaymentUseCase interface {
	GetPaymentMethodInfoByID(ctx context.Context, paymentMethodID int) (domain.PaymentMethodInfo, error)

	RazorPayCheckout(ctx context.Context, orderInfo domain.Order) (string, error)

	UpdateOrderAndPaymentDetails(ctx context.Context, paymentVerifier request.PaymentVerification) (domain.Order, error)
}
