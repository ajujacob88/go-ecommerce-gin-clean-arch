package usecase

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
)

type couponUseCase struct {
	couponRepo interfaces.CouponRepository
}

func NewCouponUseCase(couponRepo interfaces.CouponRepository) services.CouponUseCase {
	return &couponUseCase{
		couponRepo: couponRepo,
	}
}

func (c *couponUseCase) ApplyCouponToCart(ctx context.Context, userID int, couponCode string) (response.ViewCart, error) {

	//fetch the coupon with the given coupon code
	coupon, err := c.couponRepo.FetchCouponByCouponCode(ctx, couponCode)
	if err != nil {
		return response.ViewCart{}, err
	} else if coupon.ID == 0 {
		return response.ViewCart{}, fmt.Errorf("invalid coupon_id %s", couponCode)
	}

	// now check the coupon is already used by the user
	couponUsed, err := c.couponRepo.FindCouponUsedByUserID(ctx, userID, coupon.ID)
}
