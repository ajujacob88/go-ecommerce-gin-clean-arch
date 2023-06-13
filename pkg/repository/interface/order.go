package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, orderInfo domain.Order, cartItems []domain.CartItems) (domain.Order, error)
	OrderLineAndClearCart(ctx context.Context, createdOrder domain.Order, cartItems []domain.CartItems) error

	ViewOrderById(ctx context.Context, userID, orderID int) (domain.Order, error)
}
