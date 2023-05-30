//go:build wireinject
// +build wireinject

package di

import (
	http "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api"
	handler "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/config"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/db"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase"
	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(
		//database connection
		db.ConnectDatabase,

		//handlers
		handler.NewUserHandler,
		handler.NewAdminHandler,
		handler.NewProductHandler,
		handler.NewCartHandler,

		//database queries
		repository.NewUserRepository,
		repository.NewAdminRepository,
		repository.NewProductRepository,
		repository.NewOTPRepository,
		repository.NewCartRepository,

		//usecase
		usecase.NewUserUseCase,
		usecase.NewAdminUseCase,
		usecase.NewProductUseCase,
		usecase.NewOTPUseCase,
		usecase.NewCartUseCase,

		//server connection
		http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
