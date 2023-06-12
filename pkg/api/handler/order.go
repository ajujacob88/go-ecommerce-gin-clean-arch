package handler

import (
	"net/http"
	"time"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase   services.OrderUseCase
	paymentUseCase services.PaymentUseCase
	cartUseCase    services.CartUseCase
	paymentHandler PaymentHandler
}

func NewOrderHandler(orderusecase services.OrderUseCase, paymentusecase services.PaymentUseCase, cartUseCase services.CartUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase:   orderusecase,
		paymentUseCase: paymentusecase,
		cartUseCase:    cartUseCase,
		//paymentHandler: paymentHandler,
	}
}

//--------------PLACE ORDER FROM CART---------

// PlaceOrderFromCart
// @Summary Place the order from cart
// @ID place-order-from-cart
// @Description This endpoint allows a user to place the order from cart.
// @Tags Order
// @Accept json
// @Produce json
// @Param place_order_details body request.PlaceOrder true "Place Order details"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /user/cart/placeorder/ [post]
func (cr *OrderHandler) PlaceOrderFromCart(c *gin.Context) {
	var placeOrderInfo request.PlaceOrder
	if err := c.Bind(&placeOrderInfo); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "unable to read the request body", err.Error(), nil))
		return
	}

	userID, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to fetch the user ID", err.Error(), nil))
		return
	}

	// paymentMethodInfo, err := cr.paymentUseCase.GetPaymentMethodInfoByID(c.Request.Context(), placeOrderInfo.PaymentMethodID)

	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to place the order", err.Error(), nil))
	// 	return
	// }

	cartItems, err := cr.cartUseCase.FindCartItemsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to fetch the cart", err.Error(), nil))
		return
	}
	placedOrderDetails, deliveryAddress, err := cr.orderUseCase.GetOrderDetails(c.Request.Context(), userID, placeOrderInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to place the order", err.Error(), nil))
		return
	}

	//now make and save the Order
	orderInfo := domain.Order{
		UserID:              uint(userID),
		OrderDate:           time.Now(),
		PaymentMethodInfoID: uint(placeOrderInfo.PaymentMethodID),
		ShippingAddressID:   deliveryAddress.ID,
		OrderTotalPrice:     placedOrderDetails.AmountToPay,
		OrderStatusID:       2, //orderplaced
	}
	switch placeOrderInfo.PaymentMethodID {
	case 1:
		cr.OrderByCashOnDelivery(c, orderInfo, cartItems)
	case 2:
		orderInfo.OrderStatusID = 1 //order pending ... first order pending , then after razor pay verifcation, set order status to placed
		cr.paymentHandler.RazorpayCheckout(c, orderInfo, cartItems)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment method selected is invalid"})
	}

}

func (cr *OrderHandler) OrderByCashOnDelivery(c *gin.Context, orderInfo domain.Order, cartItems []domain.CartItems) {

	// save the order details
	createdOrder, err := cr.orderUseCase.SaveOrder(c.Request.Context(), orderInfo, cartItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to save the order", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(200, "succesfully placed the order", createdOrder))
}

/*
// backup before splitting, delete after splitting
func (cr *OrderHandler) PlaceOrderFromCart(c *gin.Context) {
	var placeOrderInfo request.PlaceOrder
	if err := c.Bind(&placeOrderInfo); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "unable to read the request body", err.Error(), nil))
		return
	}

	userID, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to fetch the user ID", err.Error(), nil))
		return
	}

	// paymentMethodInfo, err := cr.paymentUseCase.GetPaymentMethodInfoByID(c.Request.Context(), placeOrderInfo.PaymentMethodID)

	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to place the order", err.Error(), nil))
	// 	return
	// }

	cartItems, err := cr.cartUseCase.FindCartItemsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to fetch the cart", err.Error(), nil))
		return
	}
	placedOrderDetails, deliveryAddress, err := cr.orderUseCase.GetOrderDetails(c.Request.Context(), userID, placeOrderInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to place the order", err.Error(), nil))
		return
	}

	//now make and save the Order
	orderInfo := domain.Order{
		UserID:              uint(userID),
		OrderDate:           time.Now(),
		PaymentMethodInfoID: uint(placeOrderInfo.PaymentMethodID),
		ShippingAddressID:   deliveryAddress.ID,
		OrderTotalPrice:     placedOrderDetails.AmountToPay,
		OrderStatusID:       2,
	}

	// save the order details
	createdOrder, err := cr.orderUseCase.SaveOrder(c.Request.Context(), orderInfo, cartItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to save the order", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(200, "succesfully placed the order", createdOrder))
}
*/
