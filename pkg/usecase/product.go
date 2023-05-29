package usecase

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
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
