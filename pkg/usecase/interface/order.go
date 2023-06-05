package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
)

type OrderUseCase interface {
	GetOrderDetails(ctx context.Context, userId int, placeOrderInfo request.PlaceOrder) (response.UserOrder, domain.UserAddress, error)
}
