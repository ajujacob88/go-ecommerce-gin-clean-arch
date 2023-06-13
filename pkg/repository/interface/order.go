package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

type OrderRepository interface {
	// BeginTransaction begins a new transaction.
	BeginTransaction(ctx context.Context) error

	// Commit commits the current transaction.
	Commit(ctx context.Context) error

	// Rollback rolls back the current transaction.
	Rollback(ctx context.Context) error

	//above 3 are for the transactions inititated from usecase

	CreateOrder(ctx context.Context, orderInfo domain.Order) (domain.Order, error)
	CreatePaymentEntry(ctx context.Context, createdOrder domain.Order, paymentStatusID int) error
	OrderLineAndClearCart(ctx context.Context, createdOrder domain.Order, cartItems []domain.CartItems) error
	UpdateOrderDetails(ctx context.Context, orderID int) (domain.Order, error)

	ViewOrderById(ctx context.Context, userID, orderID int) (domain.Order, error)
}
