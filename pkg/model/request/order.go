package request

type PlaceOrder struct {
	PaymentMethodID int `json:"payment_method_id" binding:"required"`
	AddressID       int `json:"address_id" binding:"required"`
}
