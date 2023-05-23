package repository

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{DB}
}

// product category management
func (c *productDatabase) CreateCategory(ctx context.Context, newCategory string) (domain.ProductCategory, error) {
	var createdCategory domain.ProductCategory
	createCategoryQuery := `INSERT INTO product_categories(category_name)
							VALUES($1)
							RETURNING id, category_name` //By including the RETURNING clause, the INSERT statement will not only insert the new row into the table but also return the specified columns as a result. This can be useful when you need to retrieve the generated values or verify the inserted data.
	err := c.DB.Raw(createCategoryQuery, newCategory).Scan(&createdCategory).Error
	return createdCategory, err
}

//product management

func (c *productDatabase) CreateProduct(ctx context.Context, newProduct domain.Product) (domain.Product, error) {
	var createdProduct domain.Product
	productCreateQuery := `INSERT INTO products(product_category_id, name, description)
							VALUES($1,$2,$3)
							RETURNING *`
	err := c.DB.Raw(productCreateQuery, newProduct.ProductCategoryID, newProduct.Name, newProduct.Description).Scan(&createdProduct).Error
	return createdProduct, err
}
