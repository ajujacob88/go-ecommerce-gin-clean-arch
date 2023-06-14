package handler

import (
	"net/http"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	couponUseCase services.CouponUseCase
}

func NewCouponHandler(couponUseCase services.CouponUseCase) *CouponHandler {
	return &CouponHandler{
		couponUseCase: couponUseCase,
	}
}

func (cr *CouponHandler) ApplyCouponToCart(c *gin.Context) {
	var body request.ApplyCoupon
	if err := c.ShouldBindJSON(&couponCode); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "unable to fetch coupon id", err.Error(), nil))
		return
	}

	userID, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(400, "unable to fetch userid from context", err.Error(), nil))
		return
	}

	// paramID := c.Param("coupon_id")
	// couponID, err := strconv.Atoi(paramID)
	// if err != nil {
	// 	c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "unable to fetch coupon id", err.Error(), nil))
	// 	return
	// }

	cart, err := cr.couponUseCase.ApplyCouponToCart(c.Request.Context(), userID, body.CouponCode)
}
