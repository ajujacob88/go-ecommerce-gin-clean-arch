package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
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
	couponUsedQuery := `SELECT * FROM coupon_useds
						WHERE user_id = $1 AND coupon_id = $2 `
	err := c.DB.Raw(couponUsedQuery, userID, couponID).Scan(&couponUsed).Error
	if err != nil {
		return domain.CouponUsed{}, err
	}

	return couponUsed, nil
}

func (c *couponDatabase) FindCouponByCouponName(ctx context.Context, couponName string) (domain.Coupon, error) {
	var coupon domain.Coupon
	fetchCouponQuery := `SELECT * FROM coupons WHERE coupon_name = $1`
	err := c.DB.Raw(fetchCouponQuery, couponName).Scan(&coupon).Error
	if err != nil {
		return domain.Coupon{}, fmt.Errorf("faild to fetch coupon with coupon name %v", couponName)
	}

	return coupon, nil

}

func (c *couponDatabase) AddCoupon(ctx context.Context, couponDetails request.Coupon) (domain.Coupon, error) {
	//var addedCoupon domain.Coupon

	// since this is an insert query, scanning wont work..since you're performing an INSERT operation, there is no result set to scan from. The Scan function is typically used for SELECT queries to populate a struct with the retrieved data.

	// In your code, it seems that you're using the Raw method of gorm.DB to execute the INSERT query. This bypasses the automatic handling of GORM hooks and lifecycle callbacks, which include updating the created_at and updated_at fields. So its better to use the orm models and methods instead of raw query.
	// addCouponQuery := `	INSERT INTO coupons(coupon_name,coupon_code,min_order_value,discount_percent,discount_max_amount,valid_till,description)
	// 					VALUES ($1,$2,$3,$4,$5,$6,$7)`
	// err := c.DB.Exec(addCouponQuery, couponDetails.CouponName, couponDetails.CouponCode, couponDetails.MinOrderValue, couponDetails.DiscountPercent, couponDetails.DiscountMaxAmount, couponDetails.ValidTill, couponDetails.Description).Error
	// if err != nil {
	// 	return domain.Coupon{}, fmt.Errorf("failed to add coupon to the database %w", err)
	// }

	// since this is an insert query, scanning wont work..since you're performing an INSERT operation, there is no result set to scan from. The Scan function is typically used for SELECT queries to populate a struct with the retrieved data.

	newCoupon := domain.Coupon{
		CouponName:        couponDetails.CouponName,
		CouponCode:        couponDetails.CouponCode,
		MinOrderValue:     couponDetails.MinOrderValue,
		DiscountPercent:   couponDetails.DiscountPercent,
		DiscountMaxAmount: couponDetails.DiscountMaxAmount,
		ValidTill:         couponDetails.ValidTill,
		Description:       couponDetails.Description,
	}
	err := c.DB.Create(&newCoupon).Error
	if err != nil {
		return domain.Coupon{}, fmt.Errorf("failed to add coupon to the database: %w", err)
	}

	return newCoupon, nil
}

func (c *couponDatabase) UpdateCouponUsed(ctx context.Context, couponUsed domain.CouponUsed) error {
	updateCouponUsedQuery := `INSERT INTO coupon_useds ( user_id, coupon_id, used_at) VALUES ($1, $2, $3)`
	err := c.DB.Exec(updateCouponUsedQuery, couponUsed.UserID, couponUsed.CouponID, time.Now()).Error
	return err

}

func (c *couponDatabase) FetchAllCoupons(ctx context.Context, userID int) ([]response.ViewCoupons, error) {
	var allCoupons []response.ViewCoupons
	fetchQuery := `	SELECT c.coupon_name,c.coupon_code, c.min_order_value, c.discount_percent, 
					c.discount_max_amount, c.valid_till, c.Description, COALESCE(cu.used_at, NULL) AS used_at
					FROM coupons c
					LEFT JOIN coupon_useds AS cu ON cu.coupon_id = c.id AND cu.user_id = $1;`

	if err := c.DB.Raw(fetchQuery, userID).Scan(&allCoupons).Error; err != nil {
		return []response.ViewCoupons{}, err
	}
	return allCoupons, nil
}
