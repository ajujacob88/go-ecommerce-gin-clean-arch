package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB}
}

func (c *userDatabase) UserSignUp(ctx context.Context, newUser model.NewUserInfo) (model.UserDataOutput, error) {
	var userData model.UserDataOutput

	//save the user details
	UserSignUpQuery := `INSERT INTO users(first_name, last_name, email, phone, password, created_at)
						VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id,first_name, last_name, email, phone`

	err := c.DB.Raw(UserSignUpQuery, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Phone, newUser.Password).Scan(&userData).Error

	if err != nil {
		return model.UserDataOutput{}, fmt.Errorf("failed to create the user %s", newUser.FirstName)
	}

	//insert the data into userinfo table
	insertUserinfoQuery := `INSERT INTO user_infos (is_verified, is_blocked,users_id)
							VALUES ('f','f',$1);`
	err = c.DB.Exec(insertUserinfoQuery, userData.ID).Error
	if err != nil {
		return model.UserDataOutput{}, fmt.Errorf("failed to create the user(falied to copy to userinfo table) %s", newUser.FirstName)
	}

	return userData, err
}

func (c *userDatabase) FindUser(ctx context.Context, newUser model.NewUserInfo) (domain.Users, error) {
	// check email or phone match in database
	var user domain.Users
	query := `SELECT * FROM users WHERE email = ? OR phone = ?;`
	if err := c.DB.Raw(query, newUser.Email, newUser.Phone).Scan(&user).Error; err != nil {
		return user, errors.New("failed to get the user")
	}
	return user, nil
}

// OTPVerifyStatusManage method to update the verification status
func (c *userDatabase) OTPVerifyStatusManage(ctx context.Context, userEmail string, access bool) error {
	fmt.Println("access is ", access, "and id is ", userEmail)
	result := c.DB.Model(&domain.Users{}).Where("email = ?", userEmail).Update("verify_status", access).Error
	if result != nil {
		return errors.New("failed to update OTP verification status")
	}
	return nil
}

// Finds whether a email is already in the database or not and also checks whether a user is blocked or not
func (c *userDatabase) FindByEmail(ctx context.Context, email string) (domain.Users, error) {
	//var user domain.Users
	// err := c.DB.Where("Email = ?", email).Find(&user).Error
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return domain.Users{}, errors.New("invalid email")
	// 	}
	// 	return domain.Users{}, err
	// }
	var userData domain.Users
	fmt.Println("email is", email, " and users.email is")
	findUserQuery := `	SELECT users.id, users.first_name, users.last_name, users.email, users.phone, users.password, users.block_status, users.verify_status 
						FROM users 
						WHERE users.email = $1;`

	err := c.DB.Raw(findUserQuery, email).Scan(&userData).Error
	fmt.Println("error is", err)
	if userData.BlockStatus {
		return userData, errors.New("you are blocked")
	}
	fmt.Println("userdata is", userData)

	return userData, err
}
