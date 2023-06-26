package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"

	//	_ "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/verify"

	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUseCase services.PaymentUseCase
	orderUseCase   services.OrderUseCase
	cartUseCase    services.CartUseCase
}

func NewPaymentHandler(paymentUseCase services.PaymentUseCase, orderUseCase services.OrderUseCase, cartUseCase services.CartUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
		orderUseCase:   orderUseCase,
		cartUseCase:    cartUseCase,
	}
}

func (cr *PaymentHandler) RazorpayCheckout(c *gin.Context, orderInfo domain.Order) {

	fmt.Println("debug checkpoint0")
	orderInfo.OrderStatusID = 1 //order pending ... first order pending , then after razor pay verifcation, set order status to placed

	razorpayOrderID, err := cr.paymentUseCase.RazorPayCheckout(c.Request.Context(), orderInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to complete the order", err.Error(), nil))

	}
	fmt.Println("razorpayorderid is", razorpayOrderID, "and total order value is", orderInfo.OrderTotalPrice)

	// create the order as order pending, cart clearing and orderline only adter razor pay verification
	createdOrder, err := cr.orderUseCase.SaveOrderAndPayment(c.Request.Context(), orderInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to save the order", err.Error(), nil))
		return
	}

	c.HTML(200, "app.html", gin.H{ //gin.H is to fill the placeholders like this "amount": "{{.total}}"  in the html
		"total":    createdOrder.OrderTotalPrice,
		"orderid":  razorpayOrderID,
		"name":     "smartstore name",
		"email":    "smartstore@gmail.com",
		"phone_no": "7733333333",
	})

}

// Now razor pay verify/ payment success updations
// Razorpay verify
// @Summary Handling successful payment
// @ID payment-success
// @Description updating payment details upon successful payment
// @Tags Payment
// @Accept json
// @Produce json
// @Param razorpay_payment_id query string true "provide the razorpay_payment_id"
// @Param razorpay_order_id query string true "provide the razorpay_order_id"
// @Param razorpay_signature query string true "provide the razorpay_signature"
// @Param order_id query string true "provide the order_id"
// @Param order_total query string true "provide the order_total"
// @Success 202 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/payments/success/ [post]
func (cr *PaymentHandler) RazorpayVerify(c *gin.Context) {
	razorpayPaymentID := c.Query("razorpay_payment_id")
	razorpayOrderID := c.Query("razorpay_order_id")
	razorpaySignature := c.Query("razorpay_signature")
	userorderID := c.Query("order_id")
	orderID, err := strconv.Atoi(userorderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "unable to fetch order id", err.Error(), nil))
		return
	}
	razorpayOrderTotal := c.Query("order_total")
	total, err := strconv.ParseFloat(razorpayOrderTotal, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "unable to fetch order total", err.Error(), nil))
		return
	}
	userID, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "unable to fetch user id from context", err.Error(), nil))
		return
	}

	// verify the razorpay payment (by using signature as described by https://razorpay.com/docs/payments/server-integration/go/payment-gateway/build-integration)
	err = verify.VerifyRazorpayPayment(razorpayOrderID, razorpayPaymentID, razorpaySignature)
	if err != nil {
		response := response.ErrorResponse(400, "failed to verify razorpay payment", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	paymentVerifier := request.PaymentVerification{
		UserID:          userID,
		OrderID:         orderID,
		RazorpayOrderID: razorpayOrderID,
		PaymentRef:      razorpayPaymentID,
		Total:           total,
	}

	updatedOrder, err := cr.paymentUseCase.UpdateOrderAndPaymentDetails(c.Request.Context(), paymentVerifier)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to update payment details", err.Error(), nil))

	}

	//now clear the cart and create orderline
	cartItems, err := cr.cartUseCase.FindCartItemsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to fetch the cart", err.Error(), nil))
		return
	}

	err = cr.orderUseCase.OrderLineAndClearCart(c.Request.Context(), updatedOrder, cartItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to save the orderline and clear cart", err.Error(), nil))
		return
	}

	c.JSON(http.StatusAccepted, response.SuccessResponse(202, "payment success", nil))

}

/*
//no need of this code,, this is just backup before combining
// CreateRazorpayPayment
// @Summary Users can make payment using razor pay checkout
// @ID create-razorpay-payment
// @Description Users can make payment via Razorpay after placing orders
// @Tags Payment
// @Accept json
// @Produce json
// @Param order_id path string true "Order id"
// @Success 200 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/payments/razorpay/{order_id} [get]
func (cr *PaymentHandler) RazorpayCheckout(c *gin.Context) {
	paramsID := c.Param("order_id")
	orderID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "failed to read order id", err.Error(), nil))

	}
	userID, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to fetch the user ID", err.Error(), nil))
		return
	}

	order, razorpayOrderID, err := cr.paymentUseCase.RazorPayCheckout(c.Request.Context(), userID, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to complete the order", err.Error(), nil))

	}
	fmt.Println("razorpayorderid is", razorpayOrderID, "and total order value is", order.OrderTotalPrice)
	//c.HTML(200, "app.html", gin.H{ //gin.H is to fill the placeholders like this "amount": "{{.total}}"  in the html
	// 	"amount":   order.OrderTotalPrice,
	//"order_id": razorpayOrderID,
	// 	"name":     "smartstore name",
	// 	"email":    "smartstore@gmail.com",
	// 	"contact":  "7733333333",
	//})

	c.HTML(200, "app.html", gin.H{
		"total":    order.OrderTotalPrice,
		"orderid":  razorpayOrderID,
		"name":     "smartstore name",
		"email":    "smartstore@gmail.com",
		"phone_no": "7733333333",
	})

}




// Now razor pay verify/ payment success updations
// Razorpay verify
// @Summary Handling successful payment
// @ID payment-success
// @Description updating payment details upon successful payment
// @Tags Payment
// @Accept json
// @Produce json
// @Param razorpay_payment_id query string true "provide the razorpay_payment_id"
// @Param razorpay_order_id query string true "provide the razorpay_order_id"
// @Param razorpay_signature query string true "provide the razorpay_signature"
// @Param order_id query string true "provide the order_id"
// @Param order_total query string true "provide the order_total"
// @Success 202 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/payments/success/ [post]
func (cr *PaymentHandler) RazorpayVerify(c *gin.Context) {
	razorpayPaymentID := c.Query("razorpay_payment_id")
	razorpayOrderID := c.Query("razorpay_order_id")
	razorpaySignature := c.Query("razorpay_signature")
	userorderID := c.Query("order_id")
	orderID, err := strconv.Atoi(userorderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "unable to fetch order id", err.Error(), nil))
		return
	}
	razorpayOrderTotal := c.Query("order_total")
	total, err := strconv.ParseFloat(razorpayOrderTotal, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "unable to fetch user id from context", err.Error(), nil))
		return
	}
	userID, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "unable to fetch user id from context", err.Error(), nil))
		return
	}

	// verify the razorpay payment (by using signature as described by https://razorpay.com/docs/payments/server-integration/go/payment-gateway/build-integration)
	err = verify.VerifyRazorpayPayment(razorpayOrderID, razorpayPaymentID, razorpaySignature)
	if err != nil {
		response := response.ErrorResponse(400, "failed to verify razorpay payment", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	paymentVerifier := request.PaymentVerification{
		UserID:          userID,
		OrderID:         orderID,
		RazorpayOrderID: razorpayOrderID,
		PaymentRef:      razorpayPaymentID,
		Total:           total,
	}

	err = cr.paymentUseCase.UpdatePaymentDetails(c.Request.Context(), paymentVerifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to update payment details", err.Error(), nil))

	}



	c.JSON(http.StatusAccepted, response.SuccessResponse(202, "payment success", nil))

}
*/
