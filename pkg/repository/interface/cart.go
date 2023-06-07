package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
)

type CartRepository interface {
	AddToCart(ctx context.Context, productDetailsID int, userID int) (domain.CartItems, error)
	RemoveFromCart(ctx context.Context, productDetailsID int, userId int) error
	ViewCart(ctx context.Context, userId int) (response.ViewCart, error)

	CheckCartIsValidForOrder(ctx context.Context, userID int) (response.ViewCart, error)
	FindCartIDFromUserID(ctx context.Context, user_id int) (int, error)
	FindCartItemsByCartID(ctx context.Context, cart_id int) ([]domain.CartItems, error)
	FindCartItemsByUserID(ctx context.Context, user_id int) ([]domain.CartItems, error)
}
