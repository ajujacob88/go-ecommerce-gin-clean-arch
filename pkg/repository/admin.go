package repository

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
	"gorm.io/gorm"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB}
}

func (c *adminDatabase) IsSuperAdmin(ctx context.Context, adminID int) (bool, error) {
	var isSuperAdmin bool
	superAdminCheckingQuery := `SELECT is_super_admin
								FROM admins
								WHERE id = $1` //In the SQL query string, the placeholder $1 indicates the position of the first parameter that will be passed when executing the query. In this case, the value of adminID is passed as the first parameter to the Raw method.
	err := c.DB.Raw(superAdminCheckingQuery, adminID).Scan(&isSuperAdmin).Error //This line executes the SQL query using the DB.Raw method provided by the c.DB database connection. It scans the result into the isSuperAdmin variable using the &isSuperAdmin syntax. Any error that occurs during the execution is assigned to the err variable.
	return isSuperAdmin, err
}

func (c *adminDatabase) CreateAdmin(ctx context.Context, newAdminInfo model.NewAdminInfo) (domain.Admin, error) {
	var newAdmin domain.Admin
	createAdminQuery := `	INSERT INTO admins(user_name, email, phone, password,is_super_admin, is_blocked, created_at, updated_at)
						 	VALUES($1, $2, $3,$4, false, false, NOW(), NOW()) RETURNING *;`
	err := c.DB.Raw(createAdminQuery, newAdminInfo.UserName, newAdminInfo.Email, newAdminInfo.Phone, newAdminInfo.Password).Scan(&newAdmin).Error
	newAdmin.Password = "" //By setting it to an empty string before returning, the function ensures that the password is not accessible outside of the function scope.
	return newAdmin, err
}
