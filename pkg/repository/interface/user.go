package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
)

type UserRepository interface {
	FindUser(ctx context.Context, user domain.Users) (domain.Users, error)

	CreateUser(ctx context.Context, user domain.Users) (userID uint, err error)

	OTPVerifyStatusManage(ctx context.Context, userEmail string, access bool) error
}

/*
type UserRepository interface {
	FindAll(ctx context.Context) ([]domain.Users, error)
	FindByID(ctx context.Context, id uint) (domain.Users, error)
	Save(ctx context.Context, user domain.Users) (domain.Users, error)
	Delete(ctx context.Context, user domain.Users) error
}
*/
