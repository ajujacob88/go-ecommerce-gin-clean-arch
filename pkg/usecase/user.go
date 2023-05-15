package usecase

import (
	"context"
	"errors"
	"fmt"

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
	// if err := utils.CompareUsers(user, checkUser); err != nil {
	// 	return err
	// }

}

func (uc *UserUseCase) UpdateUser(ctx context.Context, user domain.User) error {
	// Perform any necessary validations or checks before updating the user

	// Update the user in the database
	err := uc.userRepository.UpdateUser(ctx, user)
	if err != nil {
		// Handle the error if the update operation fails
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
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
