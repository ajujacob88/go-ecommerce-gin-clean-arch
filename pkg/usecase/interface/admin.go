package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/common"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
)

type AdminUseCase interface {
	CreateAdmin(ctx context.Context, newAdmin request.NewAdminInfo, adminID int) (domain.Admin, error)
	AdminLogin(ctx context.Context, input request.AdminLoginInfo) (string, response.AdminDataOutput, error)

	ListAllUsers(ctx context.Context, viewUserInfo common.QueryParams) ([]domain.Users, bool, error)
	FindUserByID(ctx context.Context, userID int) (domain.Users, error)
	BlockUser(ctx context.Context, blockInfo request.BlockUser, adminID int) (domain.UserInfo, error)
	UnblockUser(ctx context.Context, userID int) (domain.UserInfo, error)

	AdminDashboard(ctx context.Context) (response.AdminDashboard, error)
}
