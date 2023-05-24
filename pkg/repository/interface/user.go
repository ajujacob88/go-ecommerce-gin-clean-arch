package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
)

type UserRepository interface {
	FindUser(ctx context.Context, user model.NewUserInfo) (domain.Users, error)

	UserSignUp(ctx context.Context, newUser model.NewUserInfo) (model.UserDataOutput, error)

	OTPVerifyStatusManage(ctx context.Context, userEmail string, access bool) error

	FindByEmail(ctx context.Context, Email string) (domain.Users, error)
}

/*
type UserRepository interface {
	FindAll(ctx context.Context) ([]domain.Users, error)
	FindByID(ctx context.Context, id uint) (domain.Users, error)
	Save(ctx context.Context, user domain.Users) (domain.Users, error)
	Delete(ctx context.Context, user domain.Users) error
}
*/
