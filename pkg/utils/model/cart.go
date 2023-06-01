package model

type CartItems struct {
	ProductItemID    uint
	Brand            string
	Name             string
	Model            string
	Quantity         uint
	ProductItemImage string
	Price            float64
	Total            float64
}

type ViewCart struct {
	CartItemsAll []CartItems
	SubTotal     float64
}

// this is for view cart from cart repo viewcart function
type CartDetails struct {
	ID       int
	SubTotal float64
}
