package usecase

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/config"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/razorpay/razorpay-go"
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

// razor pay key downloaded from https://dashboard.razorpay.com/app/website-app-settings/api-keys
// better practice is to save it in env and use viper and get it using getconfig
// const (
// 	razorpayAPIKeyID     = "rzp_test_lbL1gwQH8QK6uq"
// 	razorpayAPIKeySecret = "WXb29TEBAJ51qxt9cbYqkI8t"
// )  no need moved to env imported from env via getconfig via viper

func (c *paymentUseCase) GetPaymentMethodInfoByID(ctx context.Context, paymentMethodID int) (domain.PaymentMethodInfo, error) {
	paymentMethodInfo, err := c.paymentRepo.GetPaymentMethodInfoByID(ctx, paymentMethodID)
	return paymentMethodInfo, err
}

func (c *paymentUseCase) RazorPayCheckout(ctx context.Context, orderInfo domain.Order) (string, error) {
	//now integrate with razor pay (by using the code from razor pay)
	//client := razorpay.NewClient("<YOUR_API_KEY>", "<YOUR_API_SECRET>")
	fmt.Println("debug checkpoint2")

	razorpayAPIKeyID := config.GetConfig().RazorpayAPIKeyID
	razorpayAPIKeySecret := config.GetConfig().RazorpayAPIKeySecret

	fmt.Println("razor pay api key is", razorpayAPIKeyID, "razor pay key secret is", razorpayAPIKeySecret)

	client := razorpay.NewClient(razorpayAPIKeyID, razorpayAPIKeySecret)
	data := map[string]interface{}{
		"amount":   orderInfo.OrderTotalPrice * 100, //as per razor pay format, it includes paisa also... https://razorpay.com/docs/payments/server-integration/go/payment-gateway/build-integration/#api-sample-code
		"currency": "INR",
		"receipt":  "paymenttest_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		return "", err
	}
	razorpayOrderIDValue := body["id"]
	razorpayOrderID, ok := razorpayOrderIDValue.(string) // type assertion from interface to string. This line assigns the value of razorpayOrderIDValue to the variable razorpayOrderID, assuming that the value is of type string.
	if !ok {
		return "", fmt.Errorf("failed to assert razorpayOrderIDValue as string")
	}
	return razorpayOrderID, err
}

func (c *paymentUseCase) UpdateOrderAndPaymentDetails(ctx context.Context, paymentVerifier request.PaymentVerification) (domain.Order, error) {

	//fetch the payment details
	paymentDetails, err := c.paymentRepo.FetchPaymentDetails(ctx, paymentVerifier.OrderID)
	if err != nil {
		return domain.Order{}, err
	}
	if paymentDetails.ID == 0 {
		return domain.Order{}, fmt.Errorf("no order found")
	}
	if paymentDetails.OrderTotalPrice != paymentVerifier.Total {
		return domain.Order{}, fmt.Errorf("payment amount and order amount does not match")
	}

	updatedOrder, err := c.orderRepo.UpdateOrderDetails(ctx, int(paymentDetails.OrderID))
	if err != nil {
		return domain.Order{}, err
	}

	updatedPayment, err := c.paymentRepo.UpdatePaymentDetails(ctx, paymentVerifier)
	if err != nil {
		return domain.Order{}, err
	}

	if updatedPayment.ID == 0 {
		return domain.Order{}, fmt.Errorf("failed to update payment details")
	}
	return updatedOrder, nil
}

// no need,, just delet after razorpay handler merging
/*
func (c *paymentUseCase) RazorPayCheckout(ctx context.Context, userID, orderID int) (domain.Order, string, error) {
	// first check the payment status, if already paid, no need to proceed with payment and if not paid, then proceed with transaction.
	paymentDetails, err := c.paymentRepo.FetchPaymentDetails(ctx, orderID)
	if err != nil {
		return domain.Order{}, "", err
	} else if paymentDetails.PaymentStatusID == 2 {
		return domain.Order{}, "", fmt.Errorf("Payment already completed")

	}

	// now fetch the order details
	order, err := c.orderRepo.ViewOrderById(ctx, userID, orderID)
	if err != nil {
		return domain.Order{}, "", err
	} else if order.ID == 0 { //if no order is found
		return domain.Order{}, "", fmt.Errorf("no such order found")
	}
	//now integrate with razor pay (by using the code from razor pay)
	//client := razorpay.NewClient("<YOUR_API_KEY>", "<YOUR_API_SECRET>")

	razorpayAPIKeyID := config.GetConfig().RazorpayAPIKeyID
	razorpayAPIKeySecret := config.GetConfig().RazorpayAPIKeySecret

	fmt.Println("razor pay api key is", razorpayAPIKeyID, "razor pay key secret is", razorpayAPIKeySecret)

	client := razorpay.NewClient(razorpayAPIKeyID, razorpayAPIKeySecret)
	data := map[string]interface{}{
		"amount":   order.OrderTotalPrice * 100, //as per razor pay format, it includes paisa also... https://razorpay.com/docs/payments/server-integration/go/payment-gateway/build-integration/#api-sample-code
		"currency": "INR",
		"receipt":  "paymenttest_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		return domain.Order{}, "", err
	}
	razorpayOrderIDValue := body["id"]
	razorpayOrderID, ok := razorpayOrderIDValue.(string) // type assertion from interface to string. This line assigns the value of razorpayOrderIDValue to the variable razorpayOrderID, assuming that the value is of type string.
	if !ok {
		return domain.Order{}, "", fmt.Errorf("failed to assert razorpayOrderIDValue as string")
	}
	return order, razorpayOrderID, err
}
*/
