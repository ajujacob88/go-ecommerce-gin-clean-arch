package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/common"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase   services.OrderUseCase
	paymentUseCase services.PaymentUseCase
	cartUseCase    services.CartUseCase
	paymentHandler *PaymentHandler
}

func NewOrderHandler(orderusecase services.OrderUseCase, paymentusecase services.PaymentUseCase, cartUseCase services.CartUseCase, paymentHandler *PaymentHandler) *OrderHandler {
	return &OrderHandler{
		orderUseCase:   orderusecase,
		paymentUseCase: paymentusecase,
		cartUseCase:    cartUseCase,
		paymentHandler: paymentHandler,
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
		AppliedCouponID:     placedOrderDetails.AppliedCouponID,

		//OrderStatusID:       2, //orderplaced
	}
	switch placeOrderInfo.PaymentMethodID {
	case 1:
		cr.OrderByCashOnDelivery(c, orderInfo, cartItems)
	case 2:
		//orderInfo.OrderStatusID = 1 //order pending ... first order pending , then after razor pay verifcation, set order status to placed
		cr.paymentHandler.RazorpayCheckout(c, orderInfo)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment method selected is invalid"})
	}

}

func (cr *OrderHandler) OrderByCashOnDelivery(c *gin.Context, orderInfo domain.Order, cartItems []domain.CartItems) {

	// save the order details
	orderInfo.OrderStatusID = 2
	createdOrder, err := cr.orderUseCase.SaveOrderAndPayment(c.Request.Context(), orderInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to save the order", err.Error(), nil))
		return
	}

	err = cr.orderUseCase.OrderLineAndClearCart(c.Request.Context(), createdOrder, cartItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to save the orderline and clear cart", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(200, "succesfully placed the order", createdOrder))
}

// ------UPDATE ORDER STATUSES BY ADMIN----
// this end point allows admins to update order statuses
// @Summary Admin can update order status of any order using order_id
// @ID update-order
// @Description Endpoint for updating order statuses by admin
// @Tags Order
// @Accept json
// @Produce json
// @Param order_details body request.UpdateOrderStatuses true "Details of the order to be updated"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /admin/orders/update [patch]
func (cr *OrderHandler) UpdateOrderStatuses(c *gin.Context) {
	var orderStatuses request.UpdateOrderStatuses
	if err := c.ShouldBindJSON(&orderStatuses); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "failed to read request body", err.Error(), nil))
		return
	}
	updatedOrder, err := cr.orderUseCase.UpdateOrderStatuses(c.Request.Context(), orderStatuses)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to update the order statuses", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(200, "succesfully updated the order statuses", updatedOrder))

}

//--------------SUBMIT RETURN REQUEST---------

// Submit return request by user
// @Summary user can return a delivered order
// @ID return-request
// @Description User can request for returning the products within 10 days after delivery
// @Tags Order
// @Accept json
// @Produce json
// @Param return_req_details body request.ReturnRequest true "Return request details"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /user/orders/return/ [post]
func (cr *OrderHandler) ReturnRequest(c *gin.Context) {
	var returnReqDetails request.ReturnRequest
	if err := c.ShouldBindJSON(&returnReqDetails); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "failed to read request body", err.Error(), nil))
		return
	}

	userID, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "failed to fetch userid from context", err.Error(), nil))
		return
	}

	err = cr.orderUseCase.SubmitReturnRequest(c, userID, returnReqDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to place return request", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(200, "succesfully placed the return request", nil))

}

// ----------CANCEL ORDER--------------
// Cancel the order by user
// @Summary user can cancel an undelivered order
// @ID cancel-order
// @Description User can can cancel the order before delivery/out for delivery
// @Tags Order
// @Accept json
// @Produce json
// @Param order_id path int true "orderID of the order to be cancelled"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /user/orders/cancel/{order_id} [patch]
func (cr *OrderHandler) CancellOrder(c *gin.Context) {
	paramsID := c.Param("order_id")
	orderID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "failed to read order id from path", err.Error(), nil))
		return
	}

	userID, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "failed to fetch userid from context", err.Error(), nil))
		return
	}

	cancelledOrder, err := cr.orderUseCase.CancelOrder(c.Request.Context(), orderID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to cancel the order", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(200, "succesfully cancelled the order", cancelledOrder))

}

// ---------VIEW ALL ORDERS OF USER-----
// VIew all orders by user
// @Summary Retrieves all orders of the currently logged in user
// @ID view-all-orders
// @Description User can can view all the orders he made
// @Tags Order
// @Accept json
// @Produce json
// @Param page query int false "Enter the page no to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /user/orders/ [get]
func (cr *OrderHandler) ViewAllOrders(c *gin.Context) {
	var queryParams common.QueryParams

	queryParams.Page, _ = strconv.Atoi(c.Query("page"))
	queryParams.Limit, _ = strconv.Atoi(c.Query("limit"))

	userID, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "failed to fetch userid from context", err.Error(), nil))
		return
	}
	orders, err := cr.orderUseCase.ViewAllOrders(c.Request.Context(), userID, queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to fetch the order", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(200, "succesfully fetched the order", orders))

}
