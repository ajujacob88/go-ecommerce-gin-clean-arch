package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

type ProductUseCase interface {
	CreateCategory(ctx context.Context, newCategory string) (domain.ProductCategory, error)

	CreateProduct(ctx context.Context, newProduct domain.Product) (domain.Product, error)
}