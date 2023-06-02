package utils

import (
	"errors"

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
