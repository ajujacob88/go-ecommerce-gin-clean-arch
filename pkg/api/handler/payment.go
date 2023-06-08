package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"

	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUseCase services.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase services.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

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
// @Param razorpay_payment_id path string true "provide the razorpay_payment_id"
// @Param razorpay_order_id path string true "provide the razorpay_order_id"
// @Param order_id path string true "provide the order_id"
// @Param order_total path string true "provide the order_total"
// @Success 202 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/payments/success/ [post]
func (cr *PaymentHandler) RazorpayVerify(c *gin.Context) {
	razorpayPaymentID := c.Param("razorpay_payment_id")
	razorpayOrderID := c.Param("razorpay_order_id")
	userorderID := c.Param("order_id")
	orderID, err := strconv.Atoi(userorderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "unable to fetch order id", err.Error(), nil))
		return
	}
	razorpayOrderTotal := c.Param("order_total")
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
	c.JSON(http.StatusAccepted, response.SuccessResponse(202, "payment success", err.Error(), true))

}
