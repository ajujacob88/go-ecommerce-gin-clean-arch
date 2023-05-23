package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB: DB}
}

func (c *userDatabase) CreateUser(ctx context.Context, user domain.Users) (userID uint, err error) {
	//save the user details
	query := `INSERT INTO users(first_name, last_name, email, phone, password)
	VALUES ($1, $2, $3, $4, $5 ) RETURNING id`

	err = c.DB.Raw(query, user.FirstName, user.LastName, user.Email, user.Phone, user.Password).Scan(&userID).Error

	if err != nil {
		return 0, fmt.Errorf("failed to create the user %s", user.FirstName)
	}

	//insert the data into userinfo table
	insertUserinfoQuery := `INSERT INTO user_infos (is_verified, is_blocked,users_id)
							VALUES ('f','f',$1);`
	err = c.DB.Exec(insertUserinfoQuery, userID).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create the user(falied to copy to userinfo table) %s", user.FirstName)
	}

	fmt.Println("the user.id is", user.ID, "and userid is", userID)
	return userID, nil
}

func (c *userDatabase) FindUser(ctx context.Context, user domain.Users) (domain.Users, error) {
	// check id or email or phone match in database
	query := `SELECT * FROM users WHERE id = ? OR Email = ? OR Phone = ?`
	if err := c.DB.Raw(query, user.ID, user.Email, user.Phone).Scan(&user).Error; err != nil {
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
	var user domain.Users
	err := c.DB.Where("Email = ?", email).Find(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Users{}, errors.New("invalid email")
		}
		return domain.Users{}, err
	}

	if user.BlockStatus {
		return user, errors.New("you are blocked")
	}
	//fmt.Println("user1 is", user)

	return user, nil
}

/*  default present in repo

func (c *userDatabase) FindAll(ctx context.Context) ([]domain.Users, error) {
	var users []domain.Users
	err := c.DB.Find(&users).Error

	return users, err
}

func (c *userDatabase) FindByID(ctx context.Context, id uint) (domain.Users, error) {
	var user domain.Users
	err := c.DB.First(&user, id).Error

	return user, err
}

func (c *userDatabase) Save(ctx context.Context, user domain.Users) (domain.Users, error) {
	err := c.DB.Save(&user).Error

	return user, err
}

func (c *userDatabase) Delete(ctx context.Context, user domain.Users) error {
	err := c.DB.Delete(&user).Error

	return err
}

*/
