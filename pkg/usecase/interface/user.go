package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
)

type UserUseCase interface {
	UserSignUp(ctx context.Context, newUser model.NewUserInfo) (model.UserDataOutput, error)

	OTPVerifyStatusManage(ctx context.Context, userEmail string, access bool) error

	LoginWithEmail(ctx context.Context, user domain.Users) (domain.Users, error)

	FindByEmail(ctx context.Context, Email string) (domain.Users, error)
}

/* no need already in the code arch

type UserUseCase interface {
	FindAll(ctx context.Context) ([]domain.Users, error)
	FindByID(ctx context.Context, id uint) (domain.Users, error)
	Save(ctx context.Context, user domain.Users) (domain.Users, error)
	Delete(ctx context.Context, user domain.Users) error
}

*/
