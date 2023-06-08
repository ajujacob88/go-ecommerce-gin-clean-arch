package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
)

type PaymentRepository interface {
	GetPaymentMethodInfoByID(ctx context.Context, paymentMethodID int) (domain.PaymentMethodInfo, error)
	FetchPaymentDetails(ctx context.Context, orderID int) (domain.PaymentDetails, error)

	UpdatePaymentDetails(ctx context.Context, paymentVerifier request.PaymentVerification) (domain.PaymentDetails, error)
}
