package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB}
}

func (c *userDatabase) UserSignUp(ctx context.Context, newUser request.NewUserInfo) (response.UserDataOutput, error) {
	var userData response.UserDataOutput

	//save the user details
	UserSignUpQuery := `INSERT INTO users(first_name, last_name, email, phone, password, created_at)
						VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id,first_name, last_name, email, phone`

	err := c.DB.Raw(UserSignUpQuery, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Phone, newUser.Password).Scan(&userData).Error

	if err != nil {
		return response.UserDataOutput{}, fmt.Errorf("failed to create the user %s", newUser.FirstName)
	}

	//insert the data into userinfo table
	insertUserinfoQuery := `INSERT INTO user_infos (is_verified, is_blocked,users_id)
							VALUES ('f','f',$1);`
	err = c.DB.Exec(insertUserinfoQuery, userData.ID).Error
	if err != nil {
		return response.UserDataOutput{}, fmt.Errorf("failed to create the user(falied to copy to userinfo table) %s", newUser.FirstName)
	}

	return userData, err
}

func (c *userDatabase) FindUser(ctx context.Context, newUser request.NewUserInfo) (domain.Users, error) {
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

func (c *userDatabase) FindByEmailOrPhoneNumber(ctx context.Context, userCredentials request.UserCredentials) (domain.Users, error) {
	var userData domain.Users
	findUserQuery := `	SELECT users.id, users.first_name, users.last_name, users.email, users.phone, users.password 
						FROM users 
						WHERE users.email = $1 AND users.phone = $2;`

	err := c.DB.Raw(findUserQuery, userCredentials.Email, userCredentials.PhoneNum).Scan(&userData).Error
	if err != nil {
		return domain.Users{}, err
	} else if userData.ID == 0 {
		return domain.Users{}, errors.New("user with such email or phone number does not exist in database")
	}
	//check this condition also
	// if userData.BlockStatus {
	// 	return userData, errors.New("you are blocked")
	// }

	return userData, err
}

// used in userlogin handler
func (c *userDatabase) BlockStatus(ctx context.Context, userId uint) (bool, error) {

	blockStatusQuery := `SELECT is_blocked FROM user_infos WHERE users_id = $1;`

	var blockStatus bool

	err := c.DB.Raw(blockStatusQuery, userId).Scan(&blockStatus).Error
	return blockStatus, err
}

// ---to add address
// no need of this method i think,, check and delte the findaddresbyid method later
func (c *userDatabase) FindAddressByID(ctx context.Context, userID int) (domain.UserAddress, error) {
	var userAddress domain.UserAddress
	findAddressQuery := `SELECT * FROM user_addresses WHERE user_id = $1`
	err := c.DB.Raw(findAddressQuery, userID).Scan(&userAddress).Error
	if err != nil {
		return domain.UserAddress{}, err
	}
	return userAddress, nil
}

func (c *userDatabase) AddAddress(ctx context.Context, userAddressInput request.UserAddressInput, userID int) (domain.UserAddress, error) {
	var addedAddress domain.UserAddress

	insertAddressQuery := `	INSERT INTO user_addresses(
								user_id, house_number, street, city, district, state, pincode, landmark) 
								VALUES($1,$2,$3,$4,$5,$6, $7, $8) RETURNING *`
	err := c.DB.Raw(insertAddressQuery, userID, userAddressInput.HouseNumber, userAddressInput.Street, userAddressInput.City, userAddressInput.District, userAddressInput.State, userAddressInput.Pincode, userAddressInput.Landmark).Scan(&addedAddress).Error

	if err != nil {
		return domain.UserAddress{}, err
	}
	return addedAddress, nil

}

/*
//no need of this function since multiple address for a user will be there, so udating based on user id wont be right, so update based on address id written as the next function
func (c *userDatabase) UpdateAddressByUserID(ctx context.Context, userAddressInput model.UserAddressInput, userID int, addressID int) (domain.UserAddress, error) {
	var updatedAddress domain.UserAddress

	//	address is already there, update it
	updateAddressQuery := `	UPDATE user_addresses SET
									house_number = $1, street = $2, city = $3, district = $4, state = $5, pincode = $6, landmark = $7
									WHERE user_id = $8
									RETURNING *`
	err := c.DB.Raw(updateAddressQuery, userAddressInput.HouseNumber, userAddressInput.Street, userAddressInput.City, userAddressInput.District, userAddressInput.State, userAddressInput.Pincode, userAddressInput.Landmark, userID).Scan(&updatedAddress).Error
	if err != nil {
		return domain.UserAddress{}, err
	}
	return updatedAddress, nil
}

*/

func (c *userDatabase) UpdateAddress(ctx context.Context, userAddressInput request.UserAddressInput, userID, addressID int) (domain.UserAddress, error) {
	var updatedAddress domain.UserAddress

	//	address is already there, update it
	updateAddressQuery := `	UPDATE user_addresses SET
									house_number = $1, street = $2, city = $3, district = $4, state = $5, pincode = $6, landmark = $7
									WHERE id = $8 AND user_id = $9
									RETURNING *`
	result := c.DB.Raw(updateAddressQuery, userAddressInput.HouseNumber, userAddressInput.Street, userAddressInput.City, userAddressInput.District, userAddressInput.State, userAddressInput.Pincode, userAddressInput.Landmark, addressID, userID).Scan(&updatedAddress)
	if result.Error != nil {
		return domain.UserAddress{}, result.Error
	}
	// check the db.raw is exected succesfully like the conditions  id = $8 AND user_id = $9 met,,, then only we can pass the error.. may be this condition neednt required to check, beacuse in front end the user is actually seeing his only address.
	// Handle the case where no rows were updated
	if result.RowsAffected == 0 {
		return domain.UserAddress{}, errors.New("this addresssss is not mapped into this user")
	}
	return updatedAddress, nil
}

func (c *userDatabase) DeleteAddress(ctx context.Context, userID, addressID int) error {
	result := c.DB.Exec("DELETE FROM user_addresses WHERE id = $1 AND user_id = $2", addressID, userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("this address is not mapped into this user")
	}
	return nil
}

func (c *userDatabase) ListAddress(ctx context.Context, userID int) ([]response.ShowAddress, error) {
	var allAddress []response.ShowAddress

	listAddressQuery := `	SELECT id,house_number,street,city,district,state,pincode,landmark  
							FROM user_addresses							
							WHERE user_id = $1; `

	err := c.DB.Raw(listAddressQuery, userID).Scan(&allAddress).Error
	if err != nil {
		return []response.ShowAddress{}, err
	}
	return allAddress, nil
}

//  in the frontend, you can implement logic to retrieve and display only the addresses that belong to the currently logged-in user. This way, the user will only see their own addresses and won't have the option to select addresses created by other users. it is recommended to implement checks in both the frontend and backend to ensure data integrity and security.

func (c *userDatabase) FindAddress(ctx context.Context, userID int, addressID int) (domain.UserAddress, error) {
	var userAddress domain.UserAddress
	findAddressQuery := `SELECT * FROM user_addresses WHERE user_id = $1 AND id = $2`
	result := c.DB.Raw(findAddressQuery, userID, addressID).Scan(&userAddress)
	if result.Error != nil {
		return domain.UserAddress{}, result.Error
	}
	if result.RowsAffected == 0 {
		return domain.UserAddress{}, errors.New("invalid addressid / this address is not mapped into this user")
	}
	return userAddress, nil
}

func (c *userDatabase) ChangePassword(ctc context.Context, NewHashedPassword, MobileNum string) error {
	fmt.Println("phone nos is", MobileNum, "new hashed pw is", NewHashedPassword)

	// Remove the "+91" prefix from the phone number
	phoneNumber := strings.TrimPrefix(MobileNum, "+91")

	changePassword := c.DB.Model(&domain.Users{}).Where("phone=?", phoneNumber).UpdateColumn("password", NewHashedPassword)
	if changePassword.Error != nil {
		return fmt.Errorf("failed to update password: %v", changePassword.Error)
	}

	if changePassword.RowsAffected == 0 {
		errMsg := fmt.Sprintf("no rows updated for phone number: %s", MobileNum)
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}

	return nil
}
