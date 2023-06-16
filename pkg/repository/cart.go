package repository

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"

	"gorm.io/gorm"
)

type cartDatabase struct {
	DB *gorm.DB
}

func NewCartRepository(DB *gorm.DB) interfaces.CartRepository {
	return &cartDatabase{DB}
}

func (c *cartDatabase) AddToCart(ctx context.Context, productDetailsID int, userID int) (domain.CartItems, error) {

	// BEgin the transaction. The Begin() method is called on the database connection object to begin a new database transaction. A transaction is a sequence of database operations that are treated as a single unit of work. It allows you to perform multiple database operations atomically, meaning that either all the operations succeed and are committed, or if any operation fails, the entire transaction is rolled back, and no changes are made to the database.By starting a transaction using Begin(), you indicate that you want to group multiple database operations together and maintain their atomicity. This ensures data integrity and consistency in case of failures or concurrent access to the database.After starting the transaction, you can perform various database operations like executing queries, inserting/updating data, etc. on the transaction object (tx in the provided code). Once all the operations are completed successfully, you can commit the transaction using tx.Commit() to persist the changes in the database. If any error occurs during the transaction, you can roll back the transaction using tx.Rollback() to undo any changes made so far.It's important to note that transactions are typically used when you need to ensure the integrity and consistency of data across multiple operations. If you're performing a single operation that doesn't depend on or affect other operations, you may not need to use a transaction. here all database operations within this function are treated as a single atomic unit.

	tx := c.DB.Begin()

	//checking the user has a cart
	var cartID int
	cartCheckQuery := `SELECT id
						FROM carts
						WHERE user_id = $1
						LIMIT 1;` //LIMIT 1 ensures that only one row is returned. This is useful when you expect the query to match at most one row, and you want to retrieve that single row.

	err := tx.Raw(cartCheckQuery, userID).Scan(&cartID).Error

	if err != nil {
		tx.Rollback()
		return domain.CartItems{}, err
	}

	//If user has no cart, creating a new one
	if cartID == 0 {
		createCartQuery := `INSERT INTO carts(user_id, sub_total)
							VALUES ($1, $2)
							RETURNING id;`
		err := tx.Raw(createCartQuery, userID, 0).Scan(&cartID).Error
		if err != nil {
			tx.Rollback()
			return domain.CartItems{}, err
		}
	}

	// check if the stocks are available
	var qty_in_stock int
	productStockQuery := ` 	SELECT qty_in_stock 
						FROM product_details
						WHERE id = $1`
	err = tx.Raw(productStockQuery, productDetailsID).Scan(&qty_in_stock).Error
	if err != nil {
		tx.Rollback()
		return domain.CartItems{}, err
	}
	if qty_in_stock <= 0 {
		tx.Rollback()
		return domain.CartItems{}, fmt.Errorf("Failed to add to cart, Product out of stock/ not enough quantity available")
	}

	// check if the productDetails is already present in the cart
	var cartItem domain.CartItems
	cartItemQuery := `	SELECT id, quantity
						FROM cart_items
						WHERE cart_id = $1 
						AND product_details_id = $2
						LIMIT 1 `
	err = tx.Raw(cartItemQuery, cartID, productDetailsID).Scan(&cartItem).Error
	if err != nil {
		tx.Rollback()
		return domain.CartItems{}, err
	}

	//if item is not present in the cart
	if cartItem.ID == 0 {
		insertToCartQuery := `	INSERT INTO cart_items (cart_id, product_details_id, quantity)
								VALUES ($1, $2, 1)
								RETURNING *;`
		err := tx.Raw(insertToCartQuery, cartID, productDetailsID).Scan(&cartItem).Error
		if err != nil {
			tx.Rollback()
			return domain.CartItems{}, err
		}

	} else {
		// if the item is already present in the cart
		updateCartQuery := `	UPDATE cart_items 
								SET quantity = $1
								WHERE id = $2
								RETURNING *;`
		err := tx.Raw(updateCartQuery, cartItem.Quantity+1, cartItem.ID).Scan(&cartItem).Error
		if err != nil {
			tx.Rollback()
			return domain.CartItems{}, err
		}
	}

	// Now update the subtotal in cart table
	// product_details_id , qauntity and cart_id is known.
	// Now fetch the price from the product_details table

	var currentSubTotal, itemPrice float64
	err = tx.Raw("SELECT price FROM product_details WHERE id = $1", productDetailsID).Scan(&itemPrice).Error
	if err != nil {
		tx.Rollback()
		return domain.CartItems{}, err
	}

	// fetch the current subtotal from the cart table
	err = tx.Raw("SELECT sub_total FROM carts WHERE id = $1", cartItem.CartID).Scan(&currentSubTotal).Error
	if err != nil {
		tx.Rollback()
		return domain.CartItems{}, err
	}

	// add the price of the new product to the current subtotal and update it in the cart
	newSubTotal := currentSubTotal + itemPrice

	err = tx.Exec("UPDATE carts SET sub_total = $1 WHERE user_id = $2", newSubTotal, userID).Error
	if err != nil {
		tx.Rollback()
		return domain.CartItems{}, err
	}

	// this is for while placing the order, done by myself
	/* No need to reduce the qty_in_stock here..., need to reduce only while placing the order
	// Now reduce the qty_in_stock in product details table
	updateCartQuery := `	UPDATE product_details
								SET qty_in_stock = $1
								WHERE id = $2;`
	err = tx.Exec(updateCartQuery, qty_in_stock-1, productDetailsID).Error //qty_in_stock already retrieved in the beginnning part of this function
	if err != nil {
		tx.Rollback()
		return domain.CartItems{}, err
	}
	*/

	// Now commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return domain.CartItems{}, err
	}
	return cartItem, nil
}

//---------Remove FROM CART

func (c *cartDatabase) RemoveFromCart(ctx context.Context, productDetailsID int, userId int) error {
	tx := c.DB.Begin()

	// find the cart id from the carts table
	var cartID int
	findCarIDQuery := `	SELECT id
						FROM carts
						WHERE user_id = $1;`

	err := tx.Raw(findCarIDQuery, userId).Scan(&cartID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// now find the quantity
	var quantity int
	findQuantityQuery := `	SELECT quantity
							FROM cart_items
							WHERE cart_id = $1
							AND product_details_id = $2;`

	err = tx.Raw(findQuantityQuery, cartID, productDetailsID).Scan(&quantity).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// if quantity is 1, then delete the row
	if quantity == 0 {
		tx.Rollback()
		return fmt.Errorf("Nothing to remove from the cart")
	} else if quantity == 1 {
		deleteRowQuery := `	DELETE FROM cart_items
							WHERE cart_id = $1
							AND product_details_id = $2;`

		err = tx.Exec(deleteRowQuery, cartID, productDetailsID).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		updateRowQuery := `	UPDATE cart_items
							SET quantity = quantity - 1 
							WHERE cart_id = $1
							AND product_details_id = $2;`

		err = tx.Exec(updateRowQuery, cartID, productDetailsID).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Now fetch price from the product_details table
	var itemPrice float64
	fetchPriceQuery := `SELECT price
						FROM product_details
						WHERE id = $1;`

	err = tx.Raw(fetchPriceQuery, productDetailsID).Scan(&itemPrice).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println("debug check, item price is", itemPrice)

	//var updatedSubTotal float64
	subTotalPriceQuery := `	UPDATE carts
							SET sub_total = sub_total - $1
							WHERE id = $2;`

	err = tx.Exec(subTotalPriceQuery, itemPrice, cartID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// Now commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return err
}

//----VIEW CART

func (c *cartDatabase) ViewCart(ctx context.Context, userId int) (response.ViewCart, error) {
	tx := c.DB.Begin()
	//find the cart_id from the carts table
	var cartDetails response.CartDetails

	err := tx.Raw("SELECT id, sub_total FROM carts WHERE user_id = $1", userId).Scan(&cartDetails).Error
	if err != nil {
		tx.Rollback()
		return response.ViewCart{}, err
	}

	var cartItems []response.CartItems
	joinQuery := `	SELECT product_details_id, product_brands.brand_name,products.name,product_details.model_no,cart_items.quantity,product_details.product_details_image,product_details.price,(cart_items.quantity * product_details.price) AS total
					FROM cart_items
					JOIN product_details
					ON cart_items.product_details_id = product_details.id
					JOIN products
					ON products.id = product_details.product_id
					JOIN product_brands
					ON product_brands.id = products.brand_id
					JOIN carts
                    ON carts.id = cart_items.cart_id
					WHERE cart_items.cart_id = $1 `

	rows, err := tx.Raw(joinQuery, cartDetails.ID).Rows()
	if err != nil {
		tx.Rollback()
		return response.ViewCart{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item response.CartItems
		err := rows.Scan(&item.ProductItemID, &item.Brand, &item.Name, &item.Model, &item.Quantity, &item.ProductItemImage, &item.Price, &item.Total)
		if err != nil {
			tx.Rollback()
			return response.ViewCart{}, err
		}
		cartItems = append(cartItems, item)
	}
	// Now commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return response.ViewCart{}, err
	}
	//return cartItems, err

	var viewCart response.ViewCart
	viewCart.CartItemsAll = cartItems
	viewCart.SubTotal = cartDetails.SubTotal
	return viewCart, err

}

// check the cart is valid for placing the order
func (c *cartDatabase) CheckCartIsValidForOrder(ctx context.Context, userID int) (response.ViewCart, error) {
	userCart, err := c.ViewCart(ctx, userID)
	if err != nil {
		return response.ViewCart{}, err
	}

	// check any of the product is out of stock in the cart
	// check if the stocks are available

	var outOfStockProductsID int
	productStockQuery := ` 	SELECT DISTINCT product_details.id
							FROM product_details
							INNER JOIN cart_items ON product_details.id = cart_items.product_details_id
							INNER JOIN carts ON cart_items.cart_id = carts.id
							WHERE carts.user_id = $1 AND product_details.qty_in_stock <= 0`
	err = c.DB.Raw(productStockQuery, userID).Scan(&outOfStockProductsID).Error
	if err != nil {
		return response.ViewCart{}, err
	}
	if outOfStockProductsID != 0 {
		return response.ViewCart{}, fmt.Errorf("some products are out of stock - hence cart is not valid for placing order")
	}
	return userCart, nil

}

func (c *cartDatabase) FindCartIDFromUserID(ctx context.Context, user_id int) (int, error) {
	var cart_id int
	err := c.DB.Raw("SELECT id FROM carts WHERE user_id = ?", user_id).Scan(&cart_id).Error
	if err != nil {
		return cart_id, err
	}
	return cart_id, nil
}

func (c *cartDatabase) FindCartItemsByCartID(ctx context.Context, cart_id int) ([]domain.CartItems, error) {
	var cartItems []domain.CartItems
	err := c.DB.Raw("SELECT * FROM cart_items WHERE cart_id = ?", cart_id).Scan(&cartItems).Error
	if err != nil {
		return []domain.CartItems{}, err
	}
	return cartItems, nil
}

func (c *cartDatabase) FindCartItemsByUserID(ctx context.Context, user_id int) ([]domain.CartItems, error) {
	var cart_id int
	err := c.DB.Raw("SELECT id FROM carts WHERE user_id = ?", user_id).Scan(&cart_id).Error
	if err != nil {
		return []domain.CartItems{}, err
	}
	var cartItems []domain.CartItems
	err = c.DB.Raw("SELECT * FROM cart_items WHERE cart_id = ?", cart_id).Scan(&cartItems).Error
	if err != nil {
		return []domain.CartItems{}, err
	}
	return cartItems, nil
}

func (c *cartDatabase) FindCartByUserID(ctx context.Context, userID int) (domain.Carts, error) {
	var cart domain.Carts
	err := c.DB.Raw(`SELECT * FROM carts WHERE user_id = ?`, userID).Scan(&cart).Error
	if err != nil {
		return domain.Carts{}, err
	}
	return cart, nil
}

func (c *cartDatabase) UpdateCart(ctx context.Context, cartID, couponID uint, discountAmount, totalPrice float64) error {

	updateCartQuery := `	UPDATE carts
							SET applied_coupon_id = $1, discount_amount = $2, total_price = $3
							WHERE id = $4;`

	err := c.DB.Exec(updateCartQuery, couponID, discountAmount, totalPrice, cartID).Error
	if err != nil {
		return fmt.Errorf("failed to update the cart with discount price, %w", err)
	}
	return nil

}
