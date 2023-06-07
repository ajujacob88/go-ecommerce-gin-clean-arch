package usecase

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
)

type paymentUseCase struct {
	paymentRepo interfaces.PaymentRepository
	orderRepo   interfaces.OrderRepository
}

func NewPaymentUseCase(paymentRepo interfaces.PaymentRepository, orderRepo interfaces.OrderRepository) services.PaymentUseCase {
	return &paymentUseCase{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

func (c *paymentUseCase) GetPaymentMethodInfoByID(ctx context.Context, paymentMethodID int) (domain.PaymentMethodInfo, error) {
	paymentMethodInfo, err := c.paymentRepo.GetPaymentMethodInfoByID(ctx, paymentMethodID)
	return paymentMethodInfo, err
}

func (c *paymentUseCase) RazorPayCheckout(ctx context.Context, userID, orderID int) (domain.Order, string, error) {
	// first check the payment status, if already paid, no need to proceed with payment and if not paid, then proceed with transaction.
	paymentDetails, err := c.paymentRepo.FetchPaymentDetails(ctx, orderID)
	if err != nil {
		return domain.Order{}, "", err
	}
	if paymentDetails.PaymentStatusID == 2 {
		return domain.Order{}, "", fmt.Errorf("Payment already completed")

	}

	// now fetch the order details
	order, err := c.orderRepo.ViewOrderById(ctx, userID, orderID)
}
