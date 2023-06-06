package repository

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"gorm.io/gorm"
)

type orderDatabase struct {
	DB *gorm.DB
}

func NewOrderRepository(DB *gorm.DB) interfaces.OrderRepository {
	return &orderDatabase{DB}
}

func (c *orderDatabase) SaveOrder(ctx context.Context, orderInfo domain.Order) (domain.Order, error) {
	tx := c.DB.Begin()
	var createdOrder domain.Order
	createOrderQuery := `	INSERT INTO orders(user_id, order_date, payment_method_info_id, shipping_address_id, order_total_price, order_status_id)
							VALUES($1,$2,$3,$4,$5,$6)
							RETURNING *;`
	err := tx.Raw(createOrderQuery, orderInfo.UserID, orderInfo.OrderDate, orderInfo.PaymentMethodInfoID, orderInfo.ShippingAddressID, orderInfo.OrderTotalPrice, orderInfo.OrderStatusID).Scan(createdOrder).Error
	if err != nil {
		tx.Rollback()
		return domain.Order{}, err
	}

	//update cart_items table
	updateCartItemsQuery := `DELETE FROM cart_items WHERE cart_id = (SELECT id FROM carts WHERE user_id = $1);` //subquery written inside a query
	err = tx.Exec(updateCartItemsQuery, orderInfo.UserID).Error
	if err != nil {
		tx.Rollback()
		return domain.Order{}, err
	}

	//update carts table
	updateCartQuery := `DELETE FROM carts WHERE user_id = $1;`
	err = tx.Exec(updateCartQuery, orderInfo.UserID).Error
	if err != nil {
		tx.Rollback()
		return domain.Order{}, err
	}

	//create an entry in the payment_details table
	paymentEntryQuery := `	INSERT INTO payment_details(order_id, order_total_price, payment_method_info_id, payment_status_id, updated_at)
							VALUES ($1, $2, $3, 1, NOW());`
	err = tx.Exec(paymentEntryQuery, createdOrder.ID, createdOrder.OrderTotalPrice, createdOrder.PaymentMethodInfoID).Error
	if err != nil {
		tx.Rollback()
		return domain.Order{}, err
	}
}
