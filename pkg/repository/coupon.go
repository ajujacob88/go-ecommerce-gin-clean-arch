package repository

import (
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"gorm.io/gorm"
)

type couponDatabase struct {
	DB         *gorm.DB
	couponRepo interfaces.CouponRepository
}

func NewCouponRepository(DB *gorm.DB, couponRepo interfaces.CouponRepository) interfaces.CouponRepository {
	return &couponDatabase{
		DB:         DB,
		couponRepo: couponRepo,
	}
}
