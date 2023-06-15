package repository

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"gorm.io/gorm"
)

type couponDatabase struct {
	DB *gorm.DB
}

func NewCouponRepository(DB *gorm.DB) interfaces.CouponRepository {
	return &couponDatabase{
		DB: DB,
	}
}

func (c *couponDatabase) FetchCouponByCouponCode(ctx context.Context, couponCode string) (domain.Coupon, error) {
	var coupon domain.Coupon
	fetchCouponQuery := `SELECT * FROM coupons WHERE coupon_code = $1`
	err := c.DB.Raw(fetchCouponQuery, couponCode).Scan(&coupon).Error
	if err != nil {
		return domain.Coupon{}, fmt.Errorf("faild to fetch coupon with coupon code %v", couponCode)
	}

	return coupon, nil

}

func (c *couponDatabase) FindCouponUsedByUserIDAndCouponID(ctx context.Context, userID int, couponID uint) (domain.CouponUsed, error) {
	var couponUsed domain.CouponUsed
	couponUsedQuery := `SELECT * FROM coupon_used
						WHERE user_id = $1 AND coupon_id = $2 `
	err := c.DB.Raw(couponUsedQuery, userID, couponID).Scan(&couponUsed).Error
	if err != nil {
		return domain.CouponUsed{}, err
	}

	return couponUsed, nil
}

func (c *couponDatabase) FindCouponByCouponName(ctx context.Context, couponName string) (domain.Coupon, error) {
	var coupon domain.Coupon
	fetchCouponQuery := `SELECT * FROM coupons WHERE coupon_namee = $1`
	err := c.DB.Raw(fetchCouponQuery, couponName).Scan(&coupon).Error
	if err != nil {
		return domain.Coupon{}, fmt.Errorf("faild to fetch coupon with coupon name %v", couponName)
	}

	return coupon, nil

}

func (c *couponDatabase) AddCoupon(ctx context.Context, couponDetails domain.Coupon) (domain.Coupon, error) {
	var addedCoupon domain.Coupon
	addCouponQuery := `	INSERT INTO coupons(coupon_name,coupon_code,min_order_value,discount_percent,discount_max_amount,valid_till,description)
						VALUES ($1,$2,$3,$4,$5,$6,$7)`
	err := c.DB.Raw(addCouponQuery, couponDetails.CouponName, couponDetails.CouponCode, couponDetails.MinOrderValue, couponDetails.DiscountPercent, couponDetails.DiscountMaxAmount, couponDetails.ValidTill, couponDetails.Description).Scan(&addedCoupon).Error
	if err != nil {
		return domain.Coupon{}, fmt.Errorf("failed to add coupon to the database %w", err)
	}
	return addedCoupon, nil
}
