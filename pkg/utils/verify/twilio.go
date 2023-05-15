package verify

import (
	"errors"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

func TwilioSendOtp(phoneNumber string) (string, error) {
	//fmt.Println(phoneNumber, AUTHTOKEN, ACCOUNTSID, SERVICESID)

	//create a twilio client with twilio details
	password := config.GetConfig().AUTHTOKEN
	// fmt.Println("password", password)
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
		return "", err
	}

	return *resp.Sid, nil
}

func TwilioVerifyOTP(phoneNumber string, code string) error {
	//create a twilio client with twilio details
	password := config.GetConfig().AUTHTOKEN
	userName := config.GetConfig().ACCOUNTSID
	seviceSid := config.GetConfig().SERVICESID
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Password: password,
		Username: userName,
	})

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(seviceSid, params)
	//fmt.Println("resp is", resp, "and err is", err, "for debufg aju")
	if err != nil {
		//fmt.Println("otp not correct 1")
		return errors.New("verification check failed")
	} else if *resp.Status == "approved" {
		//fmt.Println("otp correct1")
		return nil
	} else {
		//fmt.Println("otp print 3")

		// return nil
		// if resp.Status != nil && *resp.Status == "approved" {
		// 	// Verification check was approved
		// 	return nil

		// Verification check failed or has a different status
		return errors.New("verification check failed")
	}
}
