package handler

import (
	"net/http"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase   services.OrderUseCase
	paymentUseCase services.PaymentUseCae
}

func NewOrderHandler(orderusecase services.OrderUseCase, paymentusecase services.PaymentUseCae) *OrderHandler {
	return &OrderHandler{
		orderUseCase:   orderusecase,
		paymentUseCase: paymentusecase,
	}
}

func (cr *OrderHandler) PlaceOrderFromCartCOD(c *gin.Context) {
	var placeOrderInfo request.PlaceOrder
	if err := c.Bind(&placeOrderInfo); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "unable to read the request body", err.Error(), nil))
		return
	}

	userId, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to fetch the user ID", err.Error(), nil))
		return
	}

	paymentMethodInfo, err := cr.paymentUseCase.GetPaymentMethodInfoByID(c.Request.Context(), userId)

}
