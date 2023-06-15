package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
)

type CouponUseCase interface {
	ApplyCouponToCart(ctx context.Context, userID int, couponCode string) (response.ViewCart, error)
	AddCoupon(ctx context.Context, couponDetails domain.Coupon) (domain.Coupon, error)
}
