package usecase

import (
	"context"
	"errors"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/config"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
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

func (c *otpUseCase) TwilioVerifyOTP(ctx context.Context, otpverify request.OTPVerify) (domain.OTPSession, error) {
	//create a twilio client with twilio details
	password := config.GetConfig().AUTHTOKEN
	userName := config.GetConfig().ACCOUNTSID
	seviceSid := config.GetConfig().SERVICESID
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Password: password,
		Username: userName,
	})

	otpsession, err := c.otpRepo.RetrieveOtpSession(ctx, otpverify)
	if err != nil {
		return otpsession, err
	}
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(otpsession.MobileNum)
	params.SetCode(otpverify.OTP)

	resp, err := client.VerifyV2.CreateVerificationCheck(seviceSid, params)

	if err != nil {
		return otpsession, errors.New("verification check failed")
	} else if *resp.Status == "approved" {
		return otpsession, nil
	} else {

		return otpsession, errors.New("verification check failed")
	}
}
