package handler

import (
	"net/http"
	"strconv"

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

// FindCategoryByID
// @Summary List product category by id
// @ID list-category-by-id
// @Description Admins can fetch categories by id
// @Tags Product Category
// @Accept json
// @Produce json
// @Param category_id path string true "provide the category id to be fetched"
// @Success 200 {object} res.Response
// @Failure 422 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /admin/categories/{category_id} [get]
func (cr *ProductHandler) FindCategoryByID(c *gin.Context) {
	paramsID := c.Param("id")
	categoryID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "failed to parse the category id", err.Error(), nil))
		return
	}
	category, err := cr.productUseCase.FindCategoryByID(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.ErrorResponse(500, "failed to fetch the category", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, res.SuccessResponse(200, "Succesfully fetched the category details", category))

}

// Update Category
// @Summary Admin  can update the product category details
// @ID update-category
// @Description Admins can update categories
// @Tags Product Category
// @Accept json
// @Produce json
// @Param category_details body domain.ProductCategory true "provide the category info to be updated"
// @Success 202 {object} res.Response
// @Failure 401 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /admin/categories/ [put]
func (cr *ProductHandler) UpdateCategory(c *gin.Context) {
	var updateCatInfo domain.ProductCategory
	if err := c.Bind(&updateCatInfo); err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "failed to read request body", err.Error(), nil))
		return
	}
	updatedCategory, err := cr.productUseCase.UpdateCategory(c.Request.Context(), updateCatInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(401, "unable to update the category", err.Error(), nil))
		return
	}
	c.JSON(http.StatusAccepted, res.SuccessResponse(202, "Succesfully updated the category", updatedCategory))

}

// Delete Category
// @Summary Admin  can delete the product category
// @ID delete-category
// @Description Admins can delete categories
// @Tags Product Category
// @Accept json
// @Produce json
// @Param category_id path string true "Enter the category id"
// @Success 202 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /admin/categories/{category_id} [delete]
func (cr *ProductHandler) DeleteCategory(c *gin.Context) {
	paramsID := c.Param("id")
	categoryID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "failed to parse the category id", err.Error(), nil))
		return
	}

	deletedCategory, err := cr.productUseCase.DeleteCategory(c.Request.Context(), categoryID)

	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "unable to delete the category, products listed with this category", err.Error(), nil))
		return
	}
	c.JSON(http.StatusAccepted, res.SuccessResponse(202, "Succesfully deleted the category", deletedCategory))

}

// ----------Product Brand Management

// CreateBrand
// @Summary Admin can create new product brand
// @ID create-brand
// @Description Admins can create new brands from the admin panel
// @Tags Product Brand
// @Accept json
// @Produce json
// @Param brand_name body domain.ProductBrand true "New brand name"
// @Success 201 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /admin/brands [post]
func (cr *ProductHandler) CreateBrand(c *gin.Context) {
	var newBrandDetails domain.ProductBrand
	if err := c.Bind(&newBrandDetails); err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "unable to process the request", err.Error(), nil))
		return
	}
	//  call the createbrand usecase to create a new category
	createdBrand, err := cr.productUseCase.CreateBrand(c.Request.Context(), newBrandDetails)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(400, "failed to create new brand", err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, res.SuccessResponse(201, "Brand Created Succesfully", createdBrand))

}

//------Product Management -----------

// product management
// CreateProduct
// @Summary Admin can create new product listings
// @ID create-product
// @Description Admins can create new product listings
// @Tags Products
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
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "failed to add new product", err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, res.SuccessResponse(201, "New product added succesfully", createdProduct))

}

// List All Products
// @Summary List All products
// @Description Admins and users can list all products
// @Tags Products
// @Accept json
// @Produce json
// @Param limit query int false "Number of items to retrieve per page"
// @Param page query int false "Enter the page no to display"
// cpmmenting - query query string false "Search query string"
// commenting - filter query string false "filter criteria for showing the products"
// @Param sort_by query string false "sorting criteria for showing the products"
// @Param sort_desc query bool false "sorting in descending order"
// @Success 200 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /user/products [get]
// @Router /admin/products [get]
func (cr *ProductHandler) ListAllProducts(c *gin.Context) {
	var viewProductsQueryParam model.QueryParams

	viewProductsQueryParam.Page, _ = strconv.Atoi(c.Query("page"))
	viewProductsQueryParam.Limit, _ = strconv.Atoi(c.Query("limit"))
	viewProductsQueryParam.Query = c.Query("query")
	viewProductsQueryParam.Filter = c.Query("filter")
	viewProductsQueryParam.SortBy = c.Query("sort_by")
	viewProductsQueryParam.SortDesc, _ = strconv.ParseBool(c.Query("sort_desc"))

	allProducts, err := cr.productUseCase.ListAllProducts(c.Request.Context(), viewProductsQueryParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, res.ErrorResponse(500, "failed to fetch the products", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, res.SuccessResponse(200, "Succesfully fetched all products", allProducts))

}

// Find Product By ID
// @Summary Admin and users can Fetch a specific product by product id
// @Description Admin and users can Fetch a specific product by product id
// @Tags Products
// @Accept json
// @Produce json
// @Param product_id path string true "provide the ID of the product to be fetched"
// @Success 200 {object} res.Response
// @Failure 422 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /user/products/{product_id} [get]
// @Router /admin/products/{product_id} [get]
func (cr *ProductHandler) FindProductByID(c *gin.Context) {
	paramsID := c.Param("id")
	productID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "failed to parse the product id", err.Error(), nil))
		return
	}
	product, err := cr.productUseCase.FindProductByID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.ErrorResponse(500, "Unable to find the product", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, res.SuccessResponse(200, "Succesfully fetched the product", product))
}

// Update Product
// @Summary Admin  can update the product details
// @ID update-product
// @Description Admins can update products
// @Tags Products
// @Accept json
// @Produce json
// @Param product_details body domain.Product true "provide the product info to be updated"
// @Success 202 {object} res.Response
// @Failure 401 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /admin/products/ [put]
func (cr *ProductHandler) UpdateProduct(c *gin.Context) {
	var updateProductInfo domain.Product
	if err := c.Bind(&updateProductInfo); err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "failed to read the request body", err.Error(), nil))
		return
	}

	updatedProduct, err := cr.productUseCase.UpdateProduct(c.Request.Context(), updateProductInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(401, "unable to update the product", err.Error(), nil))
		return
	}
	c.JSON(http.StatusAccepted, res.SuccessResponse(202, "Succesfully updated the product", updatedProduct))

}

// Delete Products
// @Summary Admin  can delete the products
// @ID delete-products
// @Description Admins can delete categories
// @Tags Products
// @Accept json
// @Produce json
// @Param product_id path string true "Enter the product id"
// @Success 202 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 404 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /admin/products/{product_id} [delete]
func (cr *ProductHandler) DeleteProduct(c *gin.Context) {
	paramsID := c.Param("id")
	productID, err := strconv.Atoi(paramsID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "failed to parse the product id", err.Error(), nil))
		return
	}

	_, err = cr.productUseCase.FindProductByID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusNotFound, res.ErrorResponse(404, "No products are available in this id", err.Error(), nil))
		return
	}

	err = cr.productUseCase.DeleteProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "unable to delete the products", err.Error(), nil))
		return
	}
	c.JSON(http.StatusAccepted, res.SuccessResponse(202, "Succesfully deleted the product", nil))
}

//--------------PRODUCT DETAILS---------

// AddProductDetails
// @Summary Add a product details
// @ID add-product-details
// @Description This endpoint allows an admin user to add the product details.
// @Tags Product Details
// @Accept json
// @Produce json
// @Param product_details body model.NewProductDetails true "Product details"
// @Success 201 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /admin/product-details/ [post]
func (cr *ProductHandler) AddProductDetails(c *gin.Context) {
	var NewProductDetails model.NewProductDetails
	if err := c.Bind(&NewProductDetails); err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "unable to process the request", err.Error(), nil))
		return
	}

	addedProdDetails, err := cr.productUseCase.AddProductDetails(c.Request.Context(), NewProductDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "failed to add the product details", err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, res.SuccessResponse(201, "Succesfully added the product details", addedProdDetails))

}
