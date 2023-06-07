package repository

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"gorm.io/gorm"
)

type orderDatabase struct {
	DB       *gorm.DB
	cartRepo interfaces.CartRepository
}

func NewOrderRepository(DB *gorm.DB, cartRepo interfaces.CartRepository) interfaces.OrderRepository {
	return &orderDatabase{
		DB:       DB,
		cartRepo: cartRepo,
	}
}

func (c *orderDatabase) SaveOrder(ctx context.Context, orderInfo domain.Order, cartItems []domain.CartItems) (domain.Order, error) {
	tx := c.DB.Begin()
	var createdOrder domain.Order
	createOrderQuery := `	INSERT INTO orders(user_id, order_date, payment_method_info_id, shipping_address_id, order_total_price, order_status_id)
							VALUES($1,$2,$3,$4,$5,$6)
							RETURNING *;`
	err := tx.Raw(createOrderQuery, orderInfo.UserID, orderInfo.OrderDate, orderInfo.PaymentMethodInfoID, orderInfo.ShippingAddressID, orderInfo.OrderTotalPrice, orderInfo.OrderStatusID).Scan(&createdOrder).Error
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

	//create the orderline entry
	orderLineEntryQuery := `	INSERT INTO order_lines(product_details_id, order_id, quantity, price)
								VALUES ($1, $2, $3, $4);`

	//before that fetch all the product_details_id in the cart and fetch the product details including price from the cart_items
	for i := range cartItems {
		// check if product is in stock and fetch product
		var productDetails struct {
			QtyInStock int //give the names same as that of product details table, if any mismatch, then data wont scan correctly
			Price      float64
		}

		prodctDetailFetchQuery := `	SELECT qty_in_stock, price 
									FROM product_details
									WHERE id = $1;`
		err := tx.Raw(prodctDetailFetchQuery, cartItems[i].ProductDetailsID).Scan(&productDetails).Error
		if err != nil {
			tx.Rollback()
			return domain.Order{}, err
		}

		// if product is out of stock
		if productDetails.QtyInStock < int(cartItems[i].Quantity) {
			tx.Rollback()
			return domain.Order{}, fmt.Errorf("product item out of stock for the id %v", cartItems[i].ProductDetailsID)
		}

		//now create the order line -- each items total price
		productItemTotalPrice := productDetails.Price * float64(cartItems[i].Quantity)
		err = tx.Exec(orderLineEntryQuery, cartItems[i].ProductDetailsID, createdOrder.ID, cartItems[i].Quantity, productItemTotalPrice).Error
		if err != nil {
			tx.Rollback()
			return domain.Order{}, err
		}

		// Now reduce the qty_in_stock in product details table
		reduceQuantityQuery := `	UPDATE product_details
									SET qty_in_stock = qty_in_stock - $1
									WHERE id = $2;`
		err = tx.Exec(reduceQuantityQuery, cartItems[i].Quantity, cartItems[i].ProductDetailsID).Error
		if err != nil {
			tx.Rollback()
			return domain.Order{}, err
		}

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
	tx.Commit()
	return createdOrder, nil
}

func (c *orderDatabase) ViewOrderById(ctx context.Context, userID, orderID int) (domain.Order, error) {
	var order domain.Order
	viewOrderQuery := `SELECT * FROM orders WHERE user_id = $1 AND id = $2;`
	err := c.DB.Raw(viewOrderQuery, userID, orderID).Scan(&order).Error

	if err != nil {
		return domain.Order{}, err
	}

	return order, err
}
