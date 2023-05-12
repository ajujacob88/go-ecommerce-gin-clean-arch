package utils

import (
	"errors"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

func CompareUsers(user, checkUser domain.Users) (err error) {
	if checkUser.Email == user.Email {
		err = errors.Join(err, errors.New("user already eists with this email"))

	}
	if checkUser.Phone == user.Phone {
		err = errors.Join(err, errors.New("user already exists with this phone number"))
	}
	return err
}
