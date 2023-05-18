package usecase

import (
	"context"
	"fmt"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
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
