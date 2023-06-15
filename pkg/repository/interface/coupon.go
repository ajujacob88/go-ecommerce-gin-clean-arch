package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

type CouponRepository interface {
	FetchCouponByCouponCode(ctx context.Context, couponCode string) (domain.Coupon, error)
	FindCouponUsedByUserIDAndCouponID(ctx context.Context, userID int, couponID uint) (domain.CouponUsed, error)

	FindCouponByCouponName(ctx context.Context, couponName string) (domain.Coupon, error)
	AddCoupon(ctx context.Context, couponDetails domain.Coupon) (domain.Coupon, error)
}
