package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
)

type ProductRepository interface {
	CreateCategory(ctx context.Context, newCategory string) (domain.ProductCategory, error)
	ListAllCategories(ctx context.Context) ([]domain.ProductCategory, error)
	FindCategoryByID(ctx context.Context, categoryID int) (domain.ProductCategory, error)
	UpdateCategory(ctx context.Context, updateCatInfo domain.ProductCategory) (domain.ProductCategory, error)
	DeleteCategory(ctx context.Context, categoryID int) (domain.ProductCategory, error)

	CreateProduct(ctx context.Context, newProduct domain.Product) (domain.Product, error)
	ListAllProducts(ctx context.Context, viewProductsQueryParam model.QueryParams) ([]domain.Product, error)
}
