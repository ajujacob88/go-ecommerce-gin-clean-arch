package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/config"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type adminUseCase struct {
	adminRepo interfaces.AdminRepository
}

func NewAdminUseCase(adminRepo interfaces.AdminRepository) services.AdminUseCase {
	return &adminUseCase{
		adminRepo: adminRepo,
	}
}

func (c *adminUseCase) CreateAdmin(ctx context.Context, newAdmin model.NewAdminInfo, adminID int) (domain.Admin, error) {
	isSuperAdmin, err := c.adminRepo.IsSuperAdmin(ctx, adminID)
	if err != nil {
		return domain.Admin{}, err
	}
	if !isSuperAdmin {
		return domain.Admin{}, fmt.Errorf("Only superadmin can create the new admins")
	}

	//hashing the admin password
	hash, err := bcrypt.GenerateFromPassword([]byte(newAdmin.Password), 10)
	if err != nil {
		return domain.Admin{}, err
	}
	newAdmin.Password = string(hash)
	newAdminOutput, err := c.adminRepo.CreateAdmin(ctx, newAdmin)
	return newAdminOutput, err

}
func (c *adminUseCase) AdminLogin(ctx context.Context, input model.AdminLoginInfo) (string, model.AdminDataOutput, error) {
	var adminDataInModel model.AdminDataOutput
	// Now find the admindata with the given email from the databse
	adminInfo, err := c.adminRepo.FindAdmin(ctx, input.Email)
	if err != nil {
		return "", adminDataInModel, fmt.Errorf("Error finding the admin")
	}
	if adminInfo.Email == "" {
		return "", adminDataInModel, fmt.Errorf("No such admin was found")

	}

	//Now compare and bcrypt the password
	if err := bcrypt.CompareHashAndPassword([]byte(adminInfo.Password), []byte(input.Password)); err != nil {
		return "", adminDataInModel, err
	}

	//Now check whether this admin is blocked by superadmin
	if adminInfo.IsBlocked {
		return "", adminDataInModel, fmt.Errorf("admin account is blocked")

	}

	// Now create JWT token and send it in cookie
	claims := jwt.MapClaims{
		"id":  adminInfo.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.GetJWTConfig()))

	//send back the created token

	//copying admin data from admindomain(from database table) to admin model struct
	// adminDataInModel.ID, adminDataInModel.UserName, adminDataInModel.Email, adminDataInModel.Phone, adminDataInModel.IsSuperAdmin = adminInfo.ID, adminInfo.UserName, adminInfo.Email, adminInfo.Phone, adminInfo.IsSuperAdmin   //t is a straightforward and concise method when you have a small number of fields to copy. However, it requires manually mapping each field, which can become cumbersome and error-prone if you have many fields or complex structures.
	copier.Copy(&adminDataInModel, &adminInfo) //Instead of using the above line for copying, we can use copier..  This method provides a more automated and flexible way of copying fields, especially when dealing with structs with a large number of fields or complex nested structures. The library handles field mapping based on struct tags, such as json, reducing the manual effort required.
	return tokenString, adminDataInModel, err
}

func (c *adminUseCase) ListAllUsers(ctx context.Context, viewUserInfo model.QueryParams) ([]domain.Users, bool, error) {
	users, isNoUsers, err := c.adminRepo.ListAllUsers(ctx, viewUserInfo)
	return users, isNoUsers, err
}

func (c *adminUseCase) FindUserByID(ctx context.Context, userID int) (domain.Users, error) {
	user, err := c.adminRepo.FindUserByID(ctx, userID)
	return user, err
}
