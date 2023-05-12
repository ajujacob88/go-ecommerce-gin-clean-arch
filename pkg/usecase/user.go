package usecase

import (
	"context"
	"errors"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}

}

func (c *userUseCase) Signup(ctx context.Context, user domain.Users) error {
	checkUser, err := c.userRepo.FindUser(ctx, user)
	if err != nil {
		return err
	}

	//if that user not exists then create new user
	if checkUser.ID == 0 {
		//hash the pasword
		hashPasswd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return errors.New("failed to hash the password")
		}
		user.Password = string(hashPasswd)

		_, err = c.userRepo.CreateUser(ctx, user)
		if err != nil {
			return err
		}
		return nil
	}
	return utils.CompareUsers(user, checkUser)
}

/*

func NewUserUseCase(repo interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (c *userUseCase) FindAll(ctx context.Context) ([]domain.Users, error) {
	users, err := c.userRepo.FindAll(ctx)
	return users, err
}

func (c *userUseCase) FindByID(ctx context.Context, id uint) (domain.Users, error) {
	user, err := c.userRepo.FindByID(ctx, id)
	return user, err
}

func (c *userUseCase) Save(ctx context.Context, user domain.Users) (domain.Users, error) {
	user, err := c.userRepo.Save(ctx, user)

	return user, err
}

func (c *userUseCase) Delete(ctx context.Context, user domain.Users) error {
	err := c.userRepo.Delete(ctx, user)

	return err
}

*/
