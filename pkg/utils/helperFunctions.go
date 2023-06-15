package utils

import (
	"errors"
	"math/rand"
	"time"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
)

func CompareUsers(newUser request.NewUserInfo, checkUser domain.Users) (err error) {
	if checkUser.Email == newUser.Email {
		err = errors.Join(err, errors.New("user already eists with this email"))

	}
	if checkUser.Phone == newUser.Phone {
		err = errors.Join(err, errors.New("user already exists with this phone number"))
	}
	return err
}

// generate coupon codes with the first 3 letters the same for all codes, the next 3 characters randomly generated, and the last 3 characters as generated numbers
func GenerateCouponCode() string {
	// Constants
	const prefix = "SMT" //  desired prefix
	const randomChars = 3
	const numberChars = 3

	// Pool of characters for random generation
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"

	// Generate random part of the coupon code
	rand.Seed(time.Now().UnixNano())
	randomCode := make([]byte, randomChars)
	for i := range randomCode {
		randomCode[i] = letters[rand.Intn(len(letters))]
	}

	// Generate random part of the coupon code as numbers
	numberCode := make([]byte, numberChars)
	for i := range numberCode {
		numberCode[i] = numbers[rand.Intn(len(numbers))]
	}

	// Concatenate prefix, random code, and number code
	couponCode := prefix + string(randomCode) + string(numberCode)

	return couponCode
}
