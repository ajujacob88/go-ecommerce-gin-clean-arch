package usecase

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
)

type orderUseCase struct {
	orderRepo      interfaces.OrderRepository
	paymentUseCase services.PaymentUseCase
	userRepo       interfaces.UserRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository, paymentUseCase services.PaymentUseCase, userRepo interfaces.UserRepository) services.OrderUseCase {
	return &orderUseCase{
		orderRepo:      orderRepo,
		paymentUseCase: paymentUseCase,
		userRepo:       userRepo,
	}
}

func (c *orderUseCase) GetOrderDetails(ctx context.Context, userID int, placeOrderInfo request.PlaceOrder) (response.UserOrder, error) {
	paymentMethodInfo, err := c.paymentUseCase.GetPaymentMethodInfoByID(ctx, placeOrderInfo.PaymentMethodID)
	if err != nil {
		return response.UserOrder{}, fmt.Errorf("failed to get payment method info: %w", err)
	}
	if paymentMethodInfo.BlockStatus {
		return response.UserOrder{}, fmt.Errorf("Selected payment method is blocked use another payment method")
	}

	// validate the addressid
	address, err := c.userRepo.FindAddress(ctx, userID, placeOrderInfo.AddressID)
	if err != nil {
		return response.UserOrder{}, err
	}
	//check the cart of the user is valid to place the order
	cart, err := c.cartRepo.



	if paymentMethodInfo.PaymentType == "CashOnDelivery" {

	}
}
