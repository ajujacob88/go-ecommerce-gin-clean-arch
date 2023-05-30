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
	fmt.Println("user is", newUser, "user.email is", newUser.Email)
	query := `SELECT * FROM users WHERE email = ? OR phone = ?;`
	if err := c.DB.Raw(query, newUser.Email, newUser.Phone).Scan(&user).Error; err != nil {
		fmt.Println("fialed to get user")
		return user, errors.New("failed to get the user")
	}
	return user, nil
}

// OTPVerifyStatusManage method to update the verification status
func (c *userDatabase) OTPVerifyStatusManage(ctx context.Context, otpsession domain.OTPSession) error {

	//but here the mob_num is in users table and is verified is in userinfo table,

	fmt.Println("otpsessiont is", otpsession.MobileNum, "otpsession without +91 is", otpsession.MobileNum[3:], "and", otpsession.OtpId)

	//using the DB.Model gorm model have the advantages of database switching and all, but inorder to learn the queries use db.exec or db.raw,,, also some companies wont use gorm model since it will slows down the execution and we will have litle control over it
	//var user domain.Users
	//var userInfo domain.UserInfo
	// err := c.DB.Model(&domain.Users{}).
	// 	Joins("JOIN user_infos ON users.id = user_infos.users_id").
	// 	Where("users.phone = ?", otpsession.MobileNum).
	// 	Updates(map[string]interface{}{
	// 		"user_infos.is_verified": true,
	// 	}).Error

	//anyway the above code is not working correctly error is there..., better use db.exec or db.raw

	// err := c.DB.Model(&domain.UserInfo{}).
	// 	Joins("users").
	// 	Where("users.phone = ?", "7736832773").
	// 	Update("user_infos.check", "hello").
	// 	Error

	//again c.db.raw causing error, so here i used db.exec,, MobileNum[3:] is to remove the +91 since in database no +91 is stored
	//use exec if no Return values is there and use raw if it is Select * from table

	err := c.DB.Exec(`UPDATE user_infos  SET is_verified = true WHERE users_id =  (
		SELECT id
	   FROM users
	   WHERE phone = $1 )`, otpsession.MobileNum[3:]).Error

	if err != nil {
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
	findUserQuery := `	SELECT users.id, users.first_name, users.last_name, users.email, users.phone, users.password 
						FROM users 
						WHERE users.email = $1;`

	err := c.DB.Raw(findUserQuery, email).Scan(&userData).Error
	fmt.Println("error is", err)
	// if userData.BlockStatus {
	// 	return userData, errors.New("you are blocked")
	// }
	fmt.Println("userdata is", userData)

	return userData, err
}

// used in userlogin handler
func (c *userDatabase) BlockStatus(ctx context.Context, userId uint) (bool, error) {

	blockStatusQuery := `SELECT is_blocked FROM user_infos WHERE users_id = $1;`

	var blockStatus bool

	err := c.DB.Raw(blockStatusQuery, userId).Scan(&blockStatus).Error
	return blockStatus, err
}
