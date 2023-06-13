package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, orderInfo domain.Order) (domain.Order, error)
	CreatePaymentEntry(ctx context.Context, createdOrder domain.Order, paymentStatusID int) error
	OrderLineAndClearCart(ctx context.Context, createdOrder domain.Order, cartItems []domain.CartItems) error

	ViewOrderById(ctx context.Context, userID, orderID int) (domain.Order, error)
}
