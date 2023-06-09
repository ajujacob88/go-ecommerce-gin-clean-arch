package verify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

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
		return *resp.Sid, err
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

		return errors.New("verification check failed")
	}
}

func VerifyRazorpayPayment(razorpayOrderID, razorpayPaymentID, razorpaySignature string) error {

	//move this to env lateron
	//razorpayAPIKeyID := "rzp_test_lbL1gwQH8QK6uq"
	razorpayAPIKeySecret := "WXb29TEBAJ51qxt9cbYqkI8t"

	//verify signature
	//for that first create the signature.. appears to be a representation of the desired operation rather than actual Go code.
	//generated_signature = hmac_sha256(razorpayOrderID+"|"+razorpayPaymentID, razorpayAPIKeySecret)  //this is a representation of the desired operation rather than actual Go code.
	signaturedata := razorpayOrderID + "|" + razorpayPaymentID
	h := hmac.New(sha256.New, []byte(razorpayAPIKeySecret))
	_, err := h.Write([]byte(signaturedata))
	if err != nil {
		fmt.Println("for debug check")
		return errors.New("faild to veify signature")

	}
	generated_signature := hex.EncodeToString(h.Sum(nil))

	fmt.Println("generated signature is", generated_signature)
	fmt.Println("razor pay signature is", razorpaySignature)

	// Compare the generated signature with the received signature
	if generated_signature != razorpaySignature {
		return errors.New("Razorpay signature does not match")
	}

	return nil

}
