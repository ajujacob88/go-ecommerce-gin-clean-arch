package domain

type PaymentMethodInfo struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	PaymentType    string `json:"payment_type" gorm:"not null"`
	BlockStatus    bool   `json:"block_status" gorm:"not null;default:false"`
	MaxAmountLimit uint   `json:"max_amount_limit" gorm:"not null"`
}
