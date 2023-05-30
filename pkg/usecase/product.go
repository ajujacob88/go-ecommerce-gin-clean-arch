package usecase

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
)

type productUseCase struct {
	productRepo interfaces.ProductRepository
}

func NewProductUseCase(productRepo interfaces.ProductRepository) services.ProductUseCase {
	return &productUseCase{
		productRepo: productRepo,
	}
}

// ---------Category Management
func (c *productUseCase) CreateCategory(ctx context.Context, newCategory string) (domain.ProductCategory, error) {
	createdCategory, err := c.productRepo.CreateCategory(ctx, newCategory)
	return createdCategory, err
}

func (c *productUseCase) ListAllCategories(ctx context.Context) ([]domain.ProductCategory, error) {
	allCategories, err := c.productRepo.ListAllCategories(ctx)
	return allCategories, err
}

func (c *productUseCase) FindCategoryByID(ctx context.Context, categoryID int) (domain.ProductCategory, error) {
	category, err := c.productRepo.FindCategoryByID(ctx, categoryID)
	return category, err
}

func (c *productUseCase) UpdateCategory(ctx context.Context, updateCatInfo domain.ProductCategory) (domain.ProductCategory, error) {
	updatedCategory, err := c.productRepo.UpdateCategory(ctx, updateCatInfo)
	return updatedCategory, err
}

func (c *productUseCase) DeleteCategory(ctx context.Context, categoryID int) (domain.ProductCategory, error) {
	deletedCategory, err := c.productRepo.DeleteCategory(ctx, categoryID)
	return deletedCategory, err
}

//----------Product Management

func (c *productUseCase) CreateProduct(ctx context.Context, newProduct domain.Product) (domain.Product, error) {
	createdProduct, err := c.productRepo.CreateProduct(ctx, newProduct)
	return createdProduct, err
}

func (c *productUseCase) ListAllProducts(ctx context.Context, viewProductsQueryParam model.QueryParams) ([]domain.Product, error) {
	allProducts, err := c.productRepo.ListAllProducts(ctx, viewProductsQueryParam)
	return allProducts, err
}

func (c *productUseCase) FindProductByID(ctx context.Context, productID int) (domain.Product, error) {
	product, err := c.productRepo.FindProductByID(ctx, productID)
	if product.Name == "" {
		return product, fmt.Errorf("invalid product id")
	}
	if productID == 0 {
		return domain.Product{}, fmt.Errorf("no product is found with that id")
	}
	return product, err
}

func (c *productUseCase) UpdateProduct(ctx context.Context, updateProductInfo domain.Product) (domain.Product, error) {
	updatedProduct, err := c.productRepo.UpdateProduct(ctx, updateProductInfo)
	return updatedProduct, err
}

func (c *productUseCase) DeleteProduct(ctx context.Context, productID int) error {
	err := c.productRepo.DeleteProduct(ctx, productID)
	return err
}
