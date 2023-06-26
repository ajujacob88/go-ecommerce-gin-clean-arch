package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
)

type UserUseCase interface {
	UserSignUp(ctx context.Context, newUser request.NewUserInfo) (response.UserDataOutput, error)

	OTPVerifyStatusManage(ctx context.Context, otpsession domain.OTPSession) error

	LoginWithEmail(ctx context.Context, user request.UserLoginEmail) (domain.Users, error)

	FindByEmail(ctx context.Context, Email string) (domain.Users, error)

	AddAddress(ctx context.Context, userAddressInput request.UserAddressInput, userID int) (domain.UserAddress, error)

	UpdateAddress(ctx context.Context, userAddressInput request.UserAddressInput, userID, addressID int) (domain.UserAddress, error)

	DeleteAddress(ctx context.Context, userID, addressID int) error

	ListAddress(ctx context.Context, userID int) ([]response.ShowAddress, error)

	FindByEmailOrPhoneNumber(ctx context.Context, userCredentials request.UserCredentials) (domain.Users, error)

	ChangePassword(ctx context.Context, NewHashedPassword, MobileNum string) error
}
