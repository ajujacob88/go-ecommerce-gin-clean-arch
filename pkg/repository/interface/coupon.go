package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

type CouponRepository interface {
	FetchCouponByCouponCode(ctx context.Context, couponCode string) (domain.Coupon, error)
	FindCouponUsedByUserIDAndCouponID(ctx context.Context, userID int, couponID uint) (domain.CouponUsed, error)
}
