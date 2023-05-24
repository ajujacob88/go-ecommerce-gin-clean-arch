package usecase

import (
	"context"
	"errors"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
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

func (c *userUseCase) UserSignUp(ctx context.Context, newUser model.NewUserInfo) (model.UserDataOutput, error) {
	checkUser, err := c.userRepo.FindUser(ctx, newUser)
	if err != nil {
		return model.UserDataOutput{}, err
	}

	//if that user not exists then create new user
	if checkUser.ID == 0 {
		//hash the pasword
		hashPasswd, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
		if err != nil {
			return model.UserDataOutput{}, errors.New("failed to hash the password")
		}
		newUser.Password = string(hashPasswd)

		userData, err := c.userRepo.UserSignUp(ctx, newUser)
		return userData, err
	}
	err = utils.CompareUsers(newUser, checkUser)
	return model.UserDataOutput{}, err

}

// Manage the otp verify status of users
func (uc *userUseCase) OTPVerifyStatusManage(ctx context.Context, userEmail string, access bool) error {
	err := uc.userRepo.OTPVerifyStatusManage(ctx, userEmail, access)
	if err != nil {
		// Handle any error
		// You can log the error or perform any other necessary actions
		return err
	}
	return nil
}

// user login
func (c *userUseCase) LoginWithEmail(ctx context.Context, user domain.Users) (domain.Users, error) {
	dbUser, dberr := c.userRepo.FindUser(ctx, model.NewUserInfo{})

	//check wether the user is found or not
	if dberr != nil {
		return user, dberr
	} else if dbUser.ID == 0 {
		return user, errors.New("user not exist with this details")
	}

	// check the user block_status to check wether user is blocked or not
	if dbUser.BlockStatus {
		return user, errors.New("user blocked by admin")
	}

	//check the user password with dbPassword
	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)) != nil {
		return user, errors.New("The entered password is wrong")
	}

	return dbUser, nil
}

func (c *userUseCase) FindByEmail(ctx context.Context, email string) (domain.Users, error) {
	users, err := c.userRepo.FindByEmail(ctx, email)
	//fmt.Println("user 2 is", users)
	return users, err
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
