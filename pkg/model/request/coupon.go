package request

type ApplyCoupon struct {
	CouponCode string `json:"coupon_code" binding:"required"`
}
