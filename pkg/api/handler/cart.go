package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/res"
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
// @Success 201 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 401 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /user/cart/add/{product_details_id} [post]
func (cr *CartHandler) AddToCart(c *gin.Context) {
	paramsID := c.Param("product_details_id")
	productDetailsID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "unable to parse the product details id", err.Error(), nil))
		return
	}

	userID, err := handlerutil.GetUserIdFromContext(c)
	fmt.Println("userid is ", userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse(401, "unable to fetch the user id from context", err.Error(), nil))
		return
	}

	cartItems, err := cr.cartUseCase.AddToCart(c.Request.Context(), productDetailsID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "failed to add products to the cart", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, res.SuccessResponse(201, "Successfully added product to the cart", cartItems))

}

// RemoveFromCart
// @Summary User can remove a product from the cart
// @ID remove-from-cart
// @Description User can remove product from the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param product_details_id path string true "product_details_id"
// @Success 204 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 401 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /user/cart/remove/{product_details_id} [post]
func (cr *CartHandler) RemoveFromCart(c *gin.Context) {
	paramsID := c.Param("product_details_id")
	productDetailsID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "unable to parse the product details id", err.Error(), nil))
		return
	}

	userId, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse(401, "unable to fetch the user id from context", err.Error(), nil))
		return
	}

	err = cr.cartUseCase.RemoveFromCart(c.Request.Context(), productDetailsID, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "failed to remove products from the cart", err.Error(), nil))
		return
	}

	c.JSON(http.StatusNoContent, res.SuccessResponse(204, "Successfully removed product from the cart", nil))

}
