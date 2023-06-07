package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, orderInfo domain.Order, cartItems []domain.CartItems) (domain.Order, error)
}
