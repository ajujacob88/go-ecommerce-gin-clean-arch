package repository

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
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

func (c *paymentDatabase) UpdatePaymentDetails(ctx context.Context, paymentVerifier request.PaymentVerification) (domain.PaymentDetails, error) {
	var updatedPayment domain.PaymentDetails
	updatePaymentQuery := `	UPDATE payment_details SET payment_method_info_id = 2, payment_status_id = 2, payment_ref = $1, updated_at = NOW()
							WHERE order_id = $2 RETURNING *;`

	err := c.DB.Raw(updatePaymentQuery, paymentRef, orderID).Scan(&updatedPayment).Error
	return updatedPayment, err
}
