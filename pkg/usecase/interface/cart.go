package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
)

type CartUseCase interface {
	AddToCart(ctx context.Context, productDetailsID int, userID int) (domain.CartItems, error)
	RemoveFromCart(ctx context.Context, productDetailsID int, userId int) error
	ViewCart(ctx context.Context, userId int) (model.ViewCart, error)
}
