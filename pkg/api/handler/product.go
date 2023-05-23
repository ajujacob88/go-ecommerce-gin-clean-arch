package handler

import (
	"net/http"

	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/res"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUseCase services.ProductUseCase
}

func NewProductHandler(usecase services.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: usecase,
	}
}

// Category Management

// CreateCategory
// @Summary Create new product category
// @ID create-category
// @Description Admins can create new categories from the admin panel
// @Tags product category
// @Accept json
// @Produce json
// @Param category_name body model.NewCategory true "New category name"
// @Success 201 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /admin/categories [post]
func (cr *ProductHandler) CreateCategory(c *gin.Context) {
	var category model.NewCategory
	if err := c.Bind(&category); err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "unable to process the request", err.Error(), nil))
		return
	}
	//  call the createcategory usecase to create a new category
	createdCategory, err := cr.productUseCase.CreateCategory(c.Request.Context(), category.CategoryName)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(400, "failed to create new category", err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, res.SuccessResponse(201, "Category Created Succesfully", createdCategory))

}
