package handler

import (
	"net/http"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
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

// ----------Category Management

// CreateCategory
// @Summary Create new product category
// @ID create-category
// @Description Admins can create new categories from the admin panel
// @Tags Product Category
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

// ListAllCategory
// @Summary List All product category
// @ID list-all-categories
// @Description Admins can list all categories from the admin panel
// @Tags Product Category
// @Accept json
// @Produce json
// @Success 200 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /admin/categories [get]
func (cr *ProductHandler) ListAllCategories(c *gin.Context) {
	categories, err := cr.productUseCase.ListAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.ErrorResponse(500, "failed to fetch the categories", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, res.SuccessResponse(200, "Succesfully fetched all categories", categories))

}

//------Product Management -----------

// product management
// CreateProduct
// @Summary Admin can create new product listings
// @ID create-product
// @Description Admins can create new product listings
// @Tags Product
// @Accept json
// @Produce json
// @Param new_product_details body domain.Product true "new product details"
// @Success 201 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /admin/products/ [post]
func (cr *ProductHandler) CreateProduct(c *gin.Context) {
	var newProduct domain.Product
	if err := c.Bind(&newProduct); err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "unable to process the request", err.Error(), nil))
		return
	}
	//  call the createcategory usecase to create a new category
	createdProduct, err := cr.productUseCase.CreateProduct(c.Request.Context(), newProduct)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(400, "failed to add new product", err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, res.SuccessResponse(201, "New product added succesfully", createdProduct))

}
