package domain

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	CouponName        string    `gorm:"unique" json:"coupon_name,omitempty"`
	CouponCode        string    `gorm:"unique, not null" json:"coupon_code,omitempty"`
	MinOrderValue     float64   `json:"min_order_value,omitempty"`
	DiscountPercent   float64   `json:"discount_percent,omitempty"`
	DiscountMaxAmount float64   `json:"discount_max_amount,omitempty"`
	ValidTill         time.Time `json:"valid_till"`
	Description       string    `json:"description"`
}

// this is to store the user who used the coupon
type CouponUsed struct {
	ID       uint      `json:"id" gorm:"primaryKey;not null"`
	CouponID uint      `json:"coupon_id" gorm:"not null"`
	Coupon   Coupon    `gorm:"foreignKey: CouponID" json:"-"`
	UserID   uint      `json:"user_id" gorm:"not null"`
	Users    Users     `gorm:"foreignKey: UserID" json:"-"`
	UsedAt   time.Time `json:"used_at" gorm:"not null"`
}
