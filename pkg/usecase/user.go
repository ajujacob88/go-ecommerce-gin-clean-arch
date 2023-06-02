package usecase

import (
	"context"
	"errors"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/req"

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

func (c *userUseCase) UserSignUp(ctx context.Context, newUser request.NewUserInfo) (response.UserDataOutput, error) {
	checkUser, err := c.userRepo.FindUser(ctx, newUser)
	if err != nil {
		return response.UserDataOutput{}, err
	}

	//if that user not exists then create new user
	if checkUser.ID == 0 {
		//hash the pasword
		hashPasswd, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
		if err != nil {
			return response.UserDataOutput{}, errors.New("failed to hash the password")
		}
		newUser.Password = string(hashPasswd)

		userData, err := c.userRepo.UserSignUp(ctx, newUser)
		return userData, err
	}
	err = utils.CompareUsers(newUser, checkUser)
	return response.UserDataOutput{}, err

}

// Manage the otp verify status of users
func (uc *userUseCase) OTPVerifyStatusManage(ctx context.Context, otpsession domain.OTPSession) error {
	err := uc.userRepo.OTPVerifyStatusManage(ctx, otpsession)
	return err
}

// user login
func (c *userUseCase) LoginWithEmail(ctx context.Context, user req.UserLoginEmail) (domain.Users, error) {

	//dbUser, dberr := c.userRepo.FindUser(ctx, user)
	dbUser, dberr := c.userRepo.FindByEmail(ctx, user.Email)

	//check wether the user is found or not
	if dberr != nil {
		return dbUser, dberr
	} else if dbUser.ID == 0 {
		return dbUser, errors.New("user not exist with this details")
	}

	// check the user block_status to check wether user is blocked or not
	// if dbUser.BlockStatus {
	// 	return user, errors.New("user blocked by admin")
	// }

	userId := dbUser.ID

	blockStatus, err := c.userRepo.BlockStatus(ctx, userId)
	if blockStatus {
		return dbUser, errors.New("The user is blocked.. Please contact support")
	}

	if err != nil {
		return dbUser, err
	}

	//check the user password with dbPassword
	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)) != nil {
		return dbUser, errors.New("The entered password is wrong")
	}

	return dbUser, nil
}

func (c *userUseCase) FindByEmail(ctx context.Context, email string) (domain.Users, error) {
	users, err := c.userRepo.FindByEmail(ctx, email)
	//fmt.Println("user 2 is", users)
	return users, err
}

//----user address

func (c *userUseCase) AddAddress(ctx context.Context, userAddressInput request.UserAddressInput, userID int) (domain.UserAddress, error) {

	address, err := c.userRepo.AddAddress(ctx, userAddressInput, userID)
	return address, err

}

func (c *userUseCase) UpdateAddress(ctx context.Context, userAddressInput request.UserAddressInput, userID, addressID int) (domain.UserAddress, error) {
	updatedAddress, err := c.userRepo.UpdateAddress(ctx, userAddressInput, userID, addressID)
	return updatedAddress, err
}

func (c *userUseCase) DeleteAddress(ctx context.Context, userID, addressID int) error {
	err := c.userRepo.DeleteAddress(ctx, userID, addressID)
	return err
}
