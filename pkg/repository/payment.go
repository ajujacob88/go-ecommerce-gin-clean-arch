package repository

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"gorm.io/gorm"
)

type paymentDatabase struct {
	DB *gorm.DB
}

func NewPaymentRepository(DB *gorm.DB) interfaces.PaymentRepository {
	return &paymentDatabase{DB}
}

func (c *paymentDatabase) GetPaymentMethodInfoByID(ctx context.Context, paymentMethodID int) (domain.PaymentMethodInfo, error) {
	var paymentmethodInfo domain.PaymentMethodInfo
	InfoFetchQuery := `	SELECT *
						FROM payment_method_infos
						WHERE id = $1`
	err := c.DB.Raw(InfoFetchQuery, paymentMethodID).Scan(&paymentmethodInfo).Error
	if err != nil {
		return domain.PaymentMethodInfo{}, fmt.Errorf("failed to fetch payment method infos by id %v \n%v", paymentMethodID, err.Error())
	}
	return paymentmethodInfo, nil
}

func (c *paymentDatabase) FetchPaymentDetails(ctx context.Context, orderID int) (domain.PaymentDetails, error) {
	var paymentDetails domain.PaymentDetails
	fetchPaymentDetailsQuery := `	SELECT *
									FROM payment_details
									WHERE order_id = $1`
	err := c.DB.Raw(fetchPaymentDetailsQuery, orderID).Scan(&paymentDetails).Error
	if err != nil {
		return domain.PaymentDetails{}, fmt.Errorf("failed to fetch payment details \n %v", err.Error())
	}
	return paymentDetails, nil
}
