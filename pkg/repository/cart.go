package repository

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
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

func (c *cartDatabase) ViewCart(ctx context.Context, userId int) ([]model.ViewCart, error) {
	tx := c.DB.Begin()
	//find the cart_id from the carts table
	var cartDetails model.CartDetails

	err := tx.Raw("SELECT id, sub_total FROM carts WHERE user_id = $1", userId).Scan(&cartDetails).Error
	if err != nil {
		tx.Rollback()
		return []model.ViewCart{}, err
	}

	var viewCart []model.ViewCart
	joinQuery := `	SELECT product_details_id, product_brands.brand_name,products.name,product_details.model_no,cart_items.quantity,product_details.product_details_image,product_details.price,(cart_items.quantity * product_details.price) AS subtotal,carts.sub_total
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
		return []model.ViewCart{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.ViewCart
		err := rows.Scan(&item.ProductItemID, &item.Brand, &item.Name, &item.Model, &item.Quantity, &item.ProductItemImage, &item.Price, &item.Total, &item.SubTotal)
		if err != nil {
			tx.Rollback()
			return []model.ViewCart{}, err
		}
		viewCart = append(viewCart, item)
	}
	// Now commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return []model.ViewCart{}, err
	}
	return viewCart, err

}
