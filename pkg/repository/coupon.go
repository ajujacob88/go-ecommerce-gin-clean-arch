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
