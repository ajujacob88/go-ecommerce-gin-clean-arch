package handler

import (
	"net/http"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
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

// ------------- ADD COUPON TO DB BY ADMIN--------
// AddCouponToDatabase
// @Summary API for admin to add the coupon
// @ID add-coupon-by-admin
// @Description This endpoint allows an admin to add coupon to the database.
// @Tags Admin Coupon
// @Accept json
// @Produce json
// @Param coupon_details body domain.Coupon true "Coupon Details"
// @Success 202 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/cart/applycoupon/ [patch]
func (cr CouponHandler) AddCoupon(c *gin.Context) {
	var couponDetails domain.Coupon
	if err := c.ShouldBindJSON(&couponDetails); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "unable to bind coupon", err.Error(), nil))
		return
	}

	addedCoupon, err := cr.couponUseCase.AddCoupon(c.Request.Context(), couponDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(500, "failed to add the coupon", err.Error(), nil))
		return
	}

	c.JSON(http.StatusAccepted, response.SuccessResponse(202, "Successfully added the coupon", addedCoupon))

}

//--------------APPLY COUPON TO CART BY USER---------

// ApplyCouponToCart
// @Summary Apply coupon to the cart
// @ID apply-coupon-to-cart
// @Description This endpoint allows a user to add coupon to the cart.
// @Tags Cart
// @Accept json
// @Produce json
// @Param coupon_code body request.ApplyCoupon true "Coupon Code"
// @Success 202 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/cart/applycoupon/ [patch]
func (cr *CouponHandler) ApplyCouponToCart(c *gin.Context) {

	//The PATCH method is used to partially update the resource at the given URL. HTTP method should be PUT or PATCH to indicate that you are updating an existing resource (the cart) with the provided coupon ID. If you write POST instead of PUT for the HTTP method in the code, it would indicate that you are creating a new resource with the provided coupon ID, rather than updating an existing resource.
	var body request.ApplyCoupon
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "unable to bind coupon code", err.Error(), nil))
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
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(500, "failed to apply the coupon", err.Error(), nil))
		return
	}

	c.JSON(http.StatusAccepted, response.SuccessResponse(202, "Successfully applied the coupon", cart))

}
