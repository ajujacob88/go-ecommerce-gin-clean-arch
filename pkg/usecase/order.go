package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
)

type orderUseCase struct {
	orderRepo      interfaces.OrderRepository
	paymentUseCase services.PaymentUseCase
	userRepo       interfaces.UserRepository
	cartRepo       interfaces.CartRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository, paymentUseCase services.PaymentUseCase, userRepo interfaces.UserRepository, cartRepo interfaces.CartRepository) services.OrderUseCase {
	return &orderUseCase{
		orderRepo:      orderRepo,
		paymentUseCase: paymentUseCase,
		userRepo:       userRepo,
		cartRepo:       cartRepo,
	}
}

func (c *orderUseCase) GetOrderDetails(ctx context.Context, userID int, placeOrderInfo request.PlaceOrder) (response.UserOrder, domain.UserAddress, error) {
	paymentMethodInfo, err := c.paymentUseCase.GetPaymentMethodInfoByID(ctx, placeOrderInfo.PaymentMethodID)
	if err != nil {
		return response.UserOrder{}, domain.UserAddress{}, fmt.Errorf("failed to get payment method info: %w", err)
	}
	if paymentMethodInfo.BlockStatus {
		return response.UserOrder{}, domain.UserAddress{}, fmt.Errorf("Selected payment method is blocked use another payment method")
	}

	// validate the addressid
	deliveryAddress, err := c.userRepo.FindAddress(ctx, userID, placeOrderInfo.AddressID)
	if err != nil {
		return response.UserOrder{}, domain.UserAddress{}, err
	}
	//check the cart of the user is valid to place the order
	userCart, err := c.cartRepo.CheckCartIsValidForOrder(ctx, userID)
	if err != nil {
		return response.UserOrder{}, domain.UserAddress{}, err
	}
	if userCart.SubTotal == 0 {
		return response.UserOrder{}, domain.UserAddress{}, errors.New("there is no products in the cart")
	}
	var userOrder response.UserOrder

	userOrder.AmountToPay = userCart.SubTotal
	if paymentMethodInfo.MaxAmountLimit < uint(userOrder.AmountToPay) {
		return response.UserOrder{}, domain.UserAddress{}, errors.New("the payment method selected is not applicable for this order, cart value is higher than selected payment method maximum transaction limit")
	}

	if paymentMethodInfo.PaymentType == "CashOnDelivery" {
		return userOrder, deliveryAddress, nil
	} else if paymentMethodInfo.PaymentType == "RazorPay" {
		return userOrder, deliveryAddress, nil
	} else {
		return response.UserOrder{}, domain.UserAddress{}, errors.New("the payment method selected is not available ")

	}
}
