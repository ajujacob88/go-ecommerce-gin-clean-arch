package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/common"
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
	couponRepo     interfaces.CouponRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository, paymentUseCase services.PaymentUseCase, userRepo interfaces.UserRepository, cartRepo interfaces.CartRepository, couponRepo interfaces.CouponRepository) services.OrderUseCase {
	return &orderUseCase{
		orderRepo:      orderRepo,
		paymentUseCase: paymentUseCase,
		userRepo:       userRepo,
		cartRepo:       cartRepo,
		couponRepo:     couponRepo,
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

	userOrder.AmountToPay = userCart.TotalPrice
	userOrder.AppliedCouponID = userCart.AppliedCouponID
	userOrder.AppliedCouponDiscount = userCart.DiscountAmount
	if paymentMethodInfo.MaxAmountLimit < uint(userOrder.AmountToPay) {
		return response.UserOrder{}, domain.UserAddress{}, errors.New("the payment method selected is not applicable for this order, cart value is higher than selected payment method maximum transaction limit")
	}

	if paymentMethodInfo.PaymentType == "COD" {
		return userOrder, deliveryAddress, nil
	} else if paymentMethodInfo.PaymentType == "RazorPay" {
		return userOrder, deliveryAddress, nil
	} else {
		return response.UserOrder{}, domain.UserAddress{}, errors.New("the payment method selected is not available / ")

	}
}

// save the order as pending, then after payment/cod verification change order status to order placed
func (c *orderUseCase) SaveOrderAndPayment(ctx context.Context, orderInfo domain.Order) (domain.Order, error) {

	// Begin the transaction  -- begin the transaction from usecase
	err := c.orderRepo.BeginTransaction(ctx)
	if err != nil {
		return domain.Order{}, err
	}

	// Defer the rollback function
	defer func() {
		if r := recover(); r != nil {
			c.orderRepo.Rollback(ctx) // Rollback the transaction on panic
		}
	}()

	createdOrder, err := c.orderRepo.CreateOrder(ctx, orderInfo)
	if err != nil {
		c.orderRepo.Rollback(ctx) // Rollback the transaction on error
		return domain.Order{}, err
	}

	//create an entry in the payment_details table - payment status id = 6 for cod & payment status id =1/pending for razor pay
	var paymentStatusID int
	if orderInfo.PaymentMethodInfoID == 1 {
		paymentStatusID = 6
	} else {
		paymentStatusID = 1
	}

	err = c.orderRepo.CreatePaymentEntry(ctx, createdOrder, paymentStatusID)
	if err != nil {
		c.orderRepo.Rollback(ctx) // Rollback the transaction on error
		return domain.Order{}, err
	}

	// Commit the transaction if everything is successful
	err = c.orderRepo.Commit(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	return createdOrder, err
}

func (c *orderUseCase) OrderLineAndClearCart(ctx context.Context, createdOrder domain.Order, cartItems []domain.CartItems) error {
	//actually the transactions should begin from usecase instead of repo.. so convert and do like that lateron

	err := c.orderRepo.OrderLineAndClearCart(ctx, createdOrder, cartItems)
	if err != nil {
		return err
	}

	//update the coupon status - make entry into the couponUsed table
	// updateCouponUsed := domain.CouponUsed{
	// 	UserID: createdOrder.UserID,
	// 	CouponID: createdOrder.AppliedCouponID,
	// }

	if createdOrder.AppliedCouponID != 0 {
		err := c.couponRepo.UpdateCouponUsed(ctx, domain.CouponUsed{
			UserID:   createdOrder.UserID,
			CouponID: createdOrder.AppliedCouponID,
		})
		if err != nil {
			return fmt.Errorf("faild to update couponUsed  for user \nerror:%v", err.Error())
		}
	}

	return nil
}

func (c *orderUseCase) UpdateOrderStatuses(ctx context.Context, orderStatuses request.UpdateOrderStatuses) (domain.Order, error) {
	updatedOrder, err := c.orderRepo.UpdateOrderStatuses(ctx, orderStatuses)
	if err != nil {
		return domain.Order{}, err
	}
	return updatedOrder, nil

}

func (c *orderUseCase) SubmitReturnRequest(ctx context.Context, userID int, returnReqDetails request.ReturnRequest) error {
	orderDetails, err := c.orderRepo.ViewOrderById(ctx, userID, int(returnReqDetails.OrderID))
	if err != nil {
		return err
	} else if orderDetails.ID == 0 {
		return errors.New("invalid order id")
	}

	if orderDetails.OrderStatusID != 6 || orderDetails.DeliveryStatusID != 3 {
		return fmt.Errorf("cannot return as order is undelivered")
	}

	if orderDetails.DeliveredAt.Sub(time.Now()) > time.Hour*24*15 {
		return fmt.Errorf("failed to place the return request as it is more than 10 days after which the order is delivered. Return period over")
	}

	// Begin the transaction  -- begin the transaction from usecase
	err = c.orderRepo.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	// Defer the rollback function
	defer func() {
		if r := recover(); r != nil {
			c.orderRepo.Rollback(ctx) // Rollback the transaction on panic
		}
	}()

	orderReturn := domain.OrderReturn{
		OrderID:      returnReqDetails.OrderID,
		ReturnReason: returnReqDetails.ReturnReason,
		RefundAmount: orderDetails.OrderTotalPrice,
		IsApproved:   false,
	}

	err = c.orderRepo.SaveOrderReturn(ctx, orderReturn)
	if err != nil {
		c.orderRepo.Rollback(ctx) // Rollback the transaction on error
		return err
	}

	//for Return Requested, the order_statuses id is 8, so change the order stauts id from 6 to 8
	returnRequestedStatusID := 8
	err = c.orderRepo.UpdateOrdersOrderStatus(ctx, orderDetails.ID, uint(returnRequestedStatusID))
	if err != nil {
		c.orderRepo.Rollback(ctx) // Rollback the transaction on error
		return err
	}

	// Commit the transaction if everything is successful
	err = c.orderRepo.Commit(ctx)
	if err != nil {
		return fmt.Errorf("faild to submit the return request \n error:%v", err)
	}

	log.Println("successfully submitted the order return request")
	return nil

}

func (c *orderUseCase) CancelOrder(ctx context.Context, orderID, userID int) (domain.Order, error) {
	order, err := c.orderRepo.ViewOrderById(ctx, userID, orderID)
	if err != nil {
		return domain.Order{}, err
	} else if order.ID == 0 {
		return domain.Order{}, errors.New("Invalid orderID")
	}

	if order.OrderStatusID != 1 && order.OrderStatusID != 2 && order.OrderStatusID != 3 && order.OrderStatusID != 4 {
		return domain.Order{}, errors.New("Can't cancell the order - order is delivered/out for delivery - Please return, if applicable")
	}

	// Begin the transaction  -- begin the transaction from usecase
	err = c.orderRepo.BeginTransaction(ctx)
	if err != nil {
		return domain.Order{}, err
	}

	// Defer the rollback function
	defer func() {
		if r := recover(); r != nil {
			c.orderRepo.Rollback(ctx) // Rollback the transaction on panic
		}
	}()

	// if orderstatus id is 1 or 2 or 3 or 4
	orderCancelledStatusID := 7
	err = c.orderRepo.UpdateOrdersOrderStatus(ctx, order.ID, uint(orderCancelledStatusID))
	if err != nil {
		c.orderRepo.Rollback(ctx) // Rollback the transaction on error
		return domain.Order{}, err
	}

	// now increase the product quantity in product_details table
	err = c.orderRepo.UpdateStockWhenOrderCancelled(ctx, order.ID)
	if err != nil {
		c.orderRepo.Rollback(ctx) // Rollback the transaction on error
		return domain.Order{}, err
	}

	// Commit the transaction if everything is successful
	err = c.orderRepo.Commit(ctx)
	if err != nil {
		return domain.Order{}, fmt.Errorf("faild to cancel the order \n error:%v", err)
	}

	order.OrderStatusID = 7 //this is to return order and inorder to avoid a database call, i just assigned value here
	log.Printf("successfully cancelled the order for order id %v", order.ID)
	return order, nil

}

func (c *orderUseCase) ViewAllOrders(ctx context.Context, userID int, queryParams common.QueryParams) ([]domain.Order, error) {
	orders, err := c.orderRepo.ViewAllOrders(ctx, userID, queryParams)
	return orders, err
}
