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

type CartHandler struct {
	cartUseCase services.CartUseCase
}

func NewCartHandler(usecase services.CartUseCase) *CartHandler {
	return &CartHandler{
		cartUseCase: usecase,
	}
}

// AddToCart
// @Summary User can add a productto the cart
// @ID add-to-cart
// @Description User can add product into the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param product_details_id path string true "product_details_id"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /user/cart/add/{product_details_id} [post]
func (cr *CartHandler) AddToCart(c *gin.Context) {
	paramsID := c.Param("product_details_id")
	productDetailsID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "unable to parse the product details id", err.Error(), nil))
		return
	}

	userID, err := handlerutil.GetUserIdFromContext(c)
	fmt.Println("userid is ", userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "unable to fetch the user id from context", err.Error(), nil))
		return
	}

	cartItems, err := cr.cartUseCase.AddToCart(c.Request.Context(), productDetailsID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to add products to the cart", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse(201, "Successfully added product to the cart", cartItems))

}

// RemoveFromCart
// @Summary User can remove a product from the cart
// @ID remove-from-cart
// @Description User can remove product from the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param product_details_id path string true "product_details_id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /user/cart/remove/{product_details_id} [delete]
func (cr *CartHandler) RemoveFromCart(c *gin.Context) {
	paramsID := c.Param("product_details_id")
	productDetailsID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "unable to parse the product details id", err.Error(), nil))
		return
	}

	userId, err := handlerutil.GetUserIdFromContext(c)

	fmt.Println("productdetailsid is", productDetailsID, "and userid is", userId)

	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "unable to fetch the user id from context", err.Error(), nil))
		return
	}

	err = cr.cartUseCase.RemoveFromCart(c.Request.Context(), productDetailsID, userId)
	fmt.Println("debug checkpoint2")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to remove products from the cart", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(200, "Successfully removed product from the cart", nil))

}

// ViewCart
// @Summary User can view the items in the cart and there total price
// @ID view-cart
// @Description User can view the cart and total price
// @Tags Cart
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/cart [get]
func (cr *CartHandler) ViewCart(c *gin.Context) {
	userId, err := handlerutil.GetUserIdFromContext(c)
	fmt.Println("userid is", userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "unable to fetch the user id from context", err.Error(), nil))
		return
	}

	viewCart, err := cr.cartUseCase.ViewCart(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "unable to fetch the cart details", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(200, "Successfully fetch the cart", viewCart))
}
