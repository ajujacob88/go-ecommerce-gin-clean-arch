package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

type CartUseCase interface {
	AddToCart(ctx context.Context, productDetailsID int, userID int) (domain.CartItems, error)
}