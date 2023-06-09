package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
)

type CartUseCase interface {
	AddToCart(ctx context.Context, productDetailsID int, userID int) (domain.CartItems, error)
	RemoveFromCart(ctx context.Context, productDetailsID int, userId int) (response.ViewCart, error)
	ViewCart(ctx context.Context, userId int) (response.ViewCart, error)
	FindCartItemsByUserID(ctx context.Context, user_id int) ([]domain.CartItems, error)
}
