package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/common"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
)

type AdminRepository interface {
	IsSuperAdmin(ctx context.Context, adminID int) (bool, error)
	CreateAdmin(ctx context.Context, newAdminInfo request.NewAdminInfo) (domain.Admin, error)
	FindAdmin(ctx context.Context, email string) (domain.Admin, error)

	ListAllUsers(ctx context.Context, queryParams common.QueryParams) ([]domain.Users, bool, error)
	FindUserByID(ctx context.Context, userID int) (domain.Users, error)
	BlockUser(ctx context.Context, blockInfo request.BlockUser, adminID int) (domain.UserInfo, error)
	UnblockUser(ctx context.Context, userID int) (domain.UserInfo, error)

	FetchOrdersSummaryData(ctx context.Context) (response.AdminDashboard, error)
	FetchTotalOrderedItems(ctx context.Context) (int, error)
	FetchTotalCreditedAmount(ctx context.Context) (float64, error)
	FetchUsersCount(ctx context.Context, adminDashboardData response.AdminDashboard) (response.AdminDashboard, error)

	FetchFullSalesReport(ctx context.Context, reqReportRange common.SalesReportDateRange) ([]response.SalesReport, error)
}
