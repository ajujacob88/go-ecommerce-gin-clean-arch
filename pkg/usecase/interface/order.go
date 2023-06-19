package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/common"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
)

type OrderUseCase interface {
	GetOrderDetails(ctx context.Context, userId int, placeOrderInfo request.PlaceOrder) (response.UserOrder, domain.UserAddress, error)
	SaveOrderAndPayment(ctx context.Context, orderInfo domain.Order) (domain.Order, error)
	OrderLineAndClearCart(ctx context.Context, createdOrder domain.Order, cartItems []domain.CartItems) error

	UpdateOrderStatuses(ctx context.Context, orderStatuses request.UpdateOrderStatuses) (domain.Order, error)
	SubmitReturnRequest(ctx context.Context, userID int, returnReqDetails request.ReturnRequest) error

	CancelOrder(ctx context.Context, orderID, userID int) (domain.Order, error)
	ViewAllOrders(ctx context.Context, userID int, queryParams common.QueryParams) ([]domain.Order, error)
}
