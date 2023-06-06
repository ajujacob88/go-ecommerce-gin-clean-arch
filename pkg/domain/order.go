package domain

import "time"

type Order struct {
	ID                  uint              `json:"id" gorm:"primaryKey"`
	UserID              uint              `json:"user_id"`
	Users               Users             `gorm:"foreignKey: UserID" json:"-"`
	OrderDate           time.Time         `json:"order_date"`
	PaymentMethodInfoID uint              `json:"payment_method_info_id"`
	PaymentMethodInfo   PaymentMethodInfo `gorm:"foreignKey:PaymentMethodInfoID" json:"-"`
	ShippingAddressID   uint              `json:"shipping_address_id"`
	UserAddress         UserAddress       `gorm:"foreignKey: ShippingAddressID" json:"-"`
	OrderTotalPrice     float64           `json:"order_total_price"`
	OrderStatusID       uint              `json:"order_status_id"`
	OrderStatus         OrderStatus       `gorm:"foreignKey: OrderStatusID" json:"-"`
	DeliveryStatusID    uint              `json:"delivery_status_id"`
	DeliveryStatus      DeliveryStatus    `gorm:"foreignKey: DeliveryStatusID" json:"-"`
	DeliveredAt         time.Time         `json:"delivered_at"`
}

type OrderStatus struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Status string `json:"status"`
}

type DeliveryStatus struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Status string `json:"status"`
}

type OrderLine struct {
	ID               uint           `gorm:"primaryKey"`
	ProductDetailsID uint           `json:"product_details_id"`
	ProductDetails   ProductDetails `gorm:"ForeignKey: ProductDetailsID" json:"-"`
	OrderID          uint           `json:"order_id"`
	Order            Order
}
