package interfaces

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
)

type AdminUseCase interface {
	CreateAdmin(ctx context.Context, newAdmin model.NewAdminInfo, adminID int) (domain.Admin, error)
}
