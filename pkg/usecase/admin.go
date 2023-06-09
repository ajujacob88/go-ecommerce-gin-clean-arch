package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/config"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/common"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"

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

func (c *adminUseCase) CreateAdmin(ctx context.Context, newAdmin request.NewAdminInfo, adminID int) (domain.Admin, error) {
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
func (c *adminUseCase) AdminLogin(ctx context.Context, input request.AdminLoginInfo) (string, response.AdminDataOutput, error) {
	var adminDataInModel response.AdminDataOutput
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

func (c *adminUseCase) ListAllUsers(ctx context.Context, viewUserInfo common.QueryParams) ([]domain.Users, bool, error) {
	users, isNoUsers, err := c.adminRepo.ListAllUsers(ctx, viewUserInfo)
	return users, isNoUsers, err
}

func (c *adminUseCase) FindUserByID(ctx context.Context, userID int) (domain.Users, error) {
	user, err := c.adminRepo.FindUserByID(ctx, userID)
	return user, err
}

func (c *adminUseCase) BlockUser(ctx context.Context, blockInfo request.BlockUser, adminID int) (domain.UserInfo, error) {
	blockedUser, err := c.adminRepo.BlockUser(ctx, blockInfo, adminID)
	return blockedUser, err
}

func (c *adminUseCase) UnblockUser(ctx context.Context, userID int) (domain.UserInfo, error) {
	unblockedUser, err := c.adminRepo.UnblockUser(ctx, userID)
	return unblockedUser, err
}

func (c *adminUseCase) AdminDashboard(ctx context.Context) (response.AdminDashboard, error) {
	var adminDashboardData response.AdminDashboard

	adminDashboardData, err := c.adminRepo.FetchOrdersSummaryData(ctx)
	if err != nil {
		return response.AdminDashboard{}, err
	}

	totalOrderedItems, err := c.adminRepo.FetchTotalOrderedItems(ctx)
	if err != nil {
		return response.AdminDashboard{}, err
	}
	adminDashboardData.TotalOrderedItems = totalOrderedItems

	totalCreditedAmount, err := c.adminRepo.FetchTotalCreditedAmount(ctx)
	if err != nil {
		return response.AdminDashboard{}, err
	}
	adminDashboardData.CreditedAmount = totalCreditedAmount
	adminDashboardData.PendingAmount = adminDashboardData.OrderValue - totalCreditedAmount

	adminDashboardData, err = c.adminRepo.FetchUsersCount(ctx, adminDashboardData)
	if err != nil {
		return response.AdminDashboard{}, err
	}
	return adminDashboardData, nil
}

func (c *adminUseCase) FetchFullSalesReport(ctx context.Context, reqReportRange common.SalesReportDateRange) ([]response.SalesReport, error) {
	fullSalesReport, err := c.adminRepo.FetchFullSalesReport(ctx, reqReportRange)
	if err != nil {
		return []response.SalesReport{}, err
	}
	log.Printf("successfully got sales report from %v to %v of limit %v",
		reqReportRange.StartDate, reqReportRange.EndDate, reqReportRange.Pagination.Count)

	return fullSalesReport, nil
}
