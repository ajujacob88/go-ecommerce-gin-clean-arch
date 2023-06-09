package repository

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"gorm.io/gorm"
)

type otpDatabase struct {
	DB *gorm.DB
}

func NewOTPRepository(DB *gorm.DB) interfaces.OTPRepository {
	return &otpDatabase{DB}
}

func (c otpDatabase) SaveOTP(ctx context.Context, resp string, phoneNumber string) error {
	// var otpsession domain.OTPSession
	// otpsession.OtpId = resp
	// otpsession.MobileNum = phoneNumber
	otpsession := domain.OTPSession{
		OtpId:     resp,
		MobileNum: phoneNumber,
	}

	err := c.DB.Create(&otpsession).Error
	return err
}

func (c otpDatabase) RetrieveOtpSession(ctx context.Context, otpverify request.OTPVerify) (domain.OTPSession, error) {
	var otpsession domain.OTPSession
	err := c.DB.Where("otp_id=?", otpverify.OtpId).Find(&otpsession).Error
	if err != nil {
		return otpsession, err
	}
	return otpsession, nil
}
