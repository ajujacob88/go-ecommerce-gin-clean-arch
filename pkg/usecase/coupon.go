package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
)

type couponUseCase struct {
	couponRepo interfaces.CouponRepository
	cartRepo   interfaces.CartRepository
}

func NewCouponUseCase(couponRepo interfaces.CouponRepository, cartRepo interfaces.CartRepository) services.CouponUseCase {
	return &couponUseCase{
		couponRepo: couponRepo,
		cartRepo:   cartRepo,
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
	couponUsed, err := c.couponRepo.FindCouponUsedByUserIDAndCouponID(ctx, userID, coupon.ID)
	if err != nil {
		return response.ViewCart{}, err
	} else if couponUsed.ID != 0 {
		return response.ViewCart{}, fmt.Errorf("user already applied this coupon at %v", couponUsed.UsedAt)
	}

	// now fetch the cart of the user
	cart, err := c.cartRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return response.ViewCart{}, err
	} else if cart.ID == 0 {
		return response.ViewCart{}, fmt.Errorf("there is no cart items avialable for the user with user_id %d", userID)

	}

	if cart.AppliedCouponID != 0 {
		return response.ViewCart{}, fmt.Errorf("cart have alreay applied a coupon with coupon id %d", cart.AppliedCouponID)
	}

	// validate the coupon expire date and cart price
	if time.Since(coupon.ValidTill) > 0 {
		return response.ViewCart{}, fmt.Errorf("Coupon Expired, Can't apply the coupon")
	}

	if cart.SubTotal < coupon.MinOrderValue {
		return response.ViewCart{}, fmt.Errorf("Can't apply the coupon \n Minimum order price should be %f inorder to apply this coupon", coupon.MinOrderValue)
	}

	// now calculate the discount for cart
	discountAmount := (cart.SubTotal * coupon.DiscountPercent) / 100

	// now cart total price will be
	totalPrice := cart.SubTotal - discountAmount

	//now update the cart
	err = c.cartRepo.UpdateCart(ctx, cart.ID, coupon.ID, discountAmount, totalPrice)
	if err != nil {
		return response.ViewCart{}, err
	}

	viewCart, err := c.cartRepo.ViewCart(ctx, userID)
	if err != nil {
		return response.ViewCart{}, err
	}

	//for recording the logs
	log.Printf("successfully updated the cart total with discount price %f", discountAmount)
	return viewCart, nil
}
