package request

type PlaceOrder struct {
	PaymentMethodID uint `json:"payment_method_id" binding:"required"`
	AddressID       uint `json:"address_id" binding:"required"`
}
