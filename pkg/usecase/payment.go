package usecase

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
)

type paymentUseCase struct {
	paymentRepo interfaces.PaymentRepository
}

func NewPaymentUseCase(paymentRepo interfaces.PaymentRepository) services.PaymentUseCae {
	return &paymentUseCase{
		paymentRepo: paymentRepo,
	}
}

func (c *paymentUseCase) GetPaymentMethodInfoByID(ctx context.Context, paymentMethodID int) (domain.PaymentMethodInfo, error) {
	paymentMethodInfo, err := c.paymentRepo.GetPaymentMethodInfoByID(ctx, paymentMethodID)
	return paymentMethodInfo, err
}
