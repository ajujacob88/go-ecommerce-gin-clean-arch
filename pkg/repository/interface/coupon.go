package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
)

type CouponRepository interface {
	FetchCouponByCouponCode(ctx context.Context, couponCode string) (domain.Coupon, error)
	FindCouponUsedByUserIDAndCouponID(ctx context.Context, userID int, couponID uint) (domain.CouponUsed, error)

	FindCouponByCouponName(ctx context.Context, couponName string) (domain.Coupon, error)
	AddCoupon(ctx context.Context, couponDetails request.Coupon) (domain.Coupon, error)

	UpdateCouponUsed(ctx context.Context, couponUsed domain.CouponUsed) error
}
