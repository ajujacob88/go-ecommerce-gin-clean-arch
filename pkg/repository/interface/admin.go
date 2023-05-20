package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
)

type AdminRepository interface {
	IsSuperAdmin(ctx context.Context, adminID int) (bool, error)
	CreateAdmin(ctx context.Context, newAdminInfo model.NewAdminInfo) (domain.Admin, error)
	FindAdmin(ctx context.Context, email string) (domain.Admin, error)
}
