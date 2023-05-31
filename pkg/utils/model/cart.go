package model

type ViewCart struct {
	ProductItemID    uint
	Brand            string
	Name             string
	Model            string
	Quantity         uint
	ProductItemImage string
	Price            float64
	Total            float64
	SubTotal         float64
}

// this is for view cart from cart repo viewcart function
type CartDetails struct {
	ID       int
	SubTotal float64
}
