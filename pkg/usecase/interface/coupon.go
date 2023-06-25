package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
)

type CouponUseCase interface {
	ApplyCouponToCart(ctx context.Context, userID int, couponCode string) (response.ViewCart, error)
	AddCoupon(ctx context.Context, couponDetails request.Coupon) (domain.Coupon, error)
	FetchAllCoupons(ctx context.Context, userID int) ([]response.ViewCoupons, error)
}
