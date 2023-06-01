package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/req"
)

type UserUseCase interface {
	UserSignUp(ctx context.Context, newUser model.NewUserInfo) (model.UserDataOutput, error)

	OTPVerifyStatusManage(ctx context.Context, otpsession domain.OTPSession) error

	LoginWithEmail(ctx context.Context, user req.UserLoginEmail) (domain.Users, error)

	FindByEmail(ctx context.Context, Email string) (domain.Users, error)

	AddAddress(ctx context.Context, userAddressInput model.UserAddressInput, userID int) (domain.UserAddress, error)

	UpdateAddress(ctx context.Context, userAddressInput model.UserAddressInput, userID, addressID int) (domain.UserAddress, error)

	DeleteAddress(ctx context.Context, userID, addressID int) error
}

/* no need already in the code arch

type UserUseCase interface {
	FindAll(ctx context.Context) ([]domain.Users, error)
	FindByID(ctx context.Context, id uint) (domain.Users, error)
	Save(ctx context.Context, user domain.Users) (domain.Users, error)
	Delete(ctx context.Context, user domain.Users) error
}

*/
