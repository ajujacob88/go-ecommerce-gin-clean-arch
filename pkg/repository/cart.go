package repository

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
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
		createCartQuery := `"INSERT INTO carts(user_id, sub_total)
							VALUES ($1, $2)
							RETURNING id;`
		err := tx.Raw(createCartQuery, userID, 0).Scan(&cartID).Error
		if err != nil {
			tx.Rollback()
			return domain.CartItems{}, err
		}
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

	// Now commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return domain.CartItems{}, err
	}
	return cartItem, nil
}
