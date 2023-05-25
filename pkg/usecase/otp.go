package usecase

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/config"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type otpUseCase struct {
	otpRepo interfaces.OTPRepository
}

func NewOTPUseCase(otpRepo interfaces.OTPRepository) services.OTPUseCase {
	return &otpUseCase{
		otpRepo: otpRepo,
	}
}

func (c *otpUseCase) TwilioSendOtp(ctx context.Context, phoneNumber string) (string, error) {
	//fmt.Println(phoneNumber, AUTHTOKEN, ACCOUNTSID, SERVICESID)

	//create a twilio client with twilio details
	password := config.GetConfig().AUTHTOKEN
	userName := config.GetConfig().ACCOUNTSID
	seviceSid := config.GetConfig().SERVICESID

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Password: password,
		Username: userName,
	})
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(seviceSid, params)
	if err != nil {
		return *resp.Sid, err
	}
	err = c.otpRepo.SaveOTP(ctx, *resp.Sid, phoneNumber)
	if err != nil {
		return *resp.Sid, err
	}
	return *resp.Sid, nil
}
