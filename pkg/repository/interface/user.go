package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
)

type UserRepository interface {
	FindUser(ctx context.Context, newUser request.NewUserInfo) (domain.Users, error)

	UserSignUp(ctx context.Context, newUser request.NewUserInfo) (model.UserDataOutput, error)

	OTPVerifyStatusManage(ctx context.Context, otpsession domain.OTPSession) error

	FindByEmail(ctx context.Context, Email string) (domain.Users, error)

	BlockStatus(ctx context.Context, userId uint) (bool, error)

	FindAddressByID(ctx context.Context, userID int) (domain.UserAddress, error)

	AddAddress(ctx context.Context, userAddressInput request.UserAddressInput, userID int) (domain.UserAddress, error)

	UpdateAddress(ctx context.Context, userAddressInput request.UserAddressInput, userID, addressID int) (domain.UserAddress, error)

	DeleteAddress(ctx context.Context, userID, addressID int) error
}
