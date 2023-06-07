package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
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
	c.HTML(200, "app.html", gin.H{ //gin.H is to fill the placeholders like this "amount": "{{.total}}"  in the html
		"amount":   order.OrderTotalPrice,
		"order_id": razorpayOrderID,
		"name":     "smartstore name",
		"email":    "smartstore@gmail.com",
		"contact":  "7733333333",
	})
}
