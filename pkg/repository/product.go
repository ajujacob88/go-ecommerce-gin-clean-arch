package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{DB}
}

// ---------product category management---------

func (c *productDatabase) CreateCategory(ctx context.Context, newCategory string) (domain.ProductCategory, error) {
	var createdCategory domain.ProductCategory
	createCategoryQuery := `INSERT INTO product_categories(category_name)
							VALUES($1)
							RETURNING id, category_name` //By including the RETURNING clause, the INSERT statement will not only insert the new row into the table but also return the specified columns as a result. This can be useful when you need to retrieve the generated values or verify the inserted data.
	err := c.DB.Raw(createCategoryQuery, newCategory).Scan(&createdCategory).Error
	return createdCategory, err
}

func (c *productDatabase) ListAllCategories(ctx context.Context) ([]domain.ProductCategory, error) {
	var allCategories []domain.ProductCategory
	listallCatQuery := `SELECT * FROM product_categories;`

	err := c.DB.Raw(listallCatQuery).Scan(&allCategories).Error

	return allCategories, err
}

func (c *productDatabase) FindCategoryByID(ctx context.Context, categoryID int) (domain.ProductCategory, error) {
	var category domain.ProductCategory
	findCatQuery := `SELECT * FROM product_categories WHERE id=$1`

	err := c.DB.Raw(findCatQuery, categoryID).Scan(&category).Error

	if category.ID == 0 {
		return domain.ProductCategory{}, fmt.Errorf("no category is found with that id")
	}
	return category, err
}

func (c *productDatabase) UpdateCategory(ctx context.Context, updateCatInfo domain.ProductCategory) (domain.ProductCategory, error) {
	var updatedCategory domain.ProductCategory
	updateCatQuery := `UPDATE product_categories
						SET category_name = $1
						WHERE id = $2
						RETURNING id, category_name` //In order to use the Scan method to map the returned values to the updatedCategory struct, you need to include the corresponding columns in the RETURNING clause of the SQL query.

	err := c.DB.Raw(updateCatQuery, updateCatInfo.CategoryName, updateCatInfo.ID).Scan(&updatedCategory).Error

	return updatedCategory, err

	/*
		//below code also good, but in below code only execution(updation happens, but returning of the updated cat from database is not happening.In the above code, the cat is updated by .Raw function and the updated cat is saving to updatedcat using scan)
		func (c *productDatabase) UpdateCategory(ctx context.Context, updateCatInfo domain.ProductCategory) (domain.ProductCategory, error) {
			// Construct the update query.
			updateCatQuery := `UPDATE product_categories
							   SET category_name = $1
							   WHERE id = $2`

			// Execute the update query.
			err := c.DB.Exec(updateCatQuery, updateCatInfo.CategoryName, updateCatInfo.ID).Error
			if err != nil {
				return domain.ProductCategory{}, err
			}

			// Return the updated category information.
			return updateCatInfo, nil
		}
	*/

}

func (c *productDatabase) DeleteCategory(ctx context.Context, categoryID int) (domain.ProductCategory, error) {
	var deletedCategory domain.ProductCategory
	deleteCatQuery := `DELETE FROM product_categories
						WHERE id = $1`
	err := c.DB.Raw(deleteCatQuery, categoryID).Scan(&deletedCategory).Error
	return deletedCategory, err
}

//---------product management

func (c *productDatabase) CreateProduct(ctx context.Context, newProduct domain.Product) (domain.Product, error) {
	var createdProduct domain.Product
	productCreateQuery := `INSERT INTO products(product_category_id, name, description)
							VALUES($1,$2,$3)
							RETURNING *`
	err := c.DB.Raw(productCreateQuery, newProduct.ProductCategoryID, newProduct.Name, newProduct.Description).Scan(&createdProduct).Error
	return createdProduct, err
}

func (c *productDatabase) ListAllProducts(ctx context.Context, viewProductsQueryParam model.QueryParams) ([]domain.Product, error) {

	findQuery := "SELECT * FROM products"
	params := []interface{}{}

	if viewProductsQueryParam.Query != "" && viewProductsQueryParam.Filter != "" {
		findQuery = fmt.Sprintf("%s WHERE LOWER(%s) LIKE $%d", findQuery, viewProductsQueryParam.Filter, len(params)+1)
		params = append(params, "%"+strings.ToLower(viewProductsQueryParam.Query)+"%")
		fmt.Println("params is ", params)
	}
	if viewProductsQueryParam.SortBy != "" {
		findQuery = fmt.Sprintf("%s ORDER BY %s %s", findQuery, viewProductsQueryParam.SortBy, orderByDirection(viewProductsQueryParam.SortDesc))
	}
	if viewProductsQueryParam.Limit != 0 && viewProductsQueryParam.Page != 0 {
		findQuery = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", findQuery, len(params)+1, len(params)+2)
		params = append(params, viewProductsQueryParam.Limit, (viewProductsQueryParam.Page-1)*viewProductsQueryParam.Limit)
	}

	var allProducts []domain.Product
	err := c.DB.Raw(findQuery, params...).Scan(&allProducts).Error

	return allProducts, err
}

func (c *productDatabase) FindProductByID(ctx context.Context, productID int) (domain.Product, error) {
	var product domain.Product
	findProductQuery := `SELECT * FROM products
						WHERE id = $1`

	err := c.DB.Raw(findProductQuery, productID).Scan(&product).Error
	return product, err
}

func (c *productDatabase) UpdateProduct(ctx context.Context, updateProductInfo domain.Product) (domain.Product, error) {
	var updatedProduct domain.Product
	updateProdQuery := `UPDATE products
						SET 
							product_category_id = $1,
							name = $2,
							description = $3
						WHERE id = $4
						RETURNING id,product_category_id,name,description`

	err := c.DB.Raw(updateProdQuery, updateProductInfo.ProductCategoryID, updateProductInfo.Name, updateProductInfo.Description, updateProductInfo.ID).Scan(&updatedProduct).Error
	return updatedProduct, err
}
