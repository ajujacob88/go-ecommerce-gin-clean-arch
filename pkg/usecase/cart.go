package usecase

import (
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
)

type cartUseCase struct {
	cartRepo interfaces.CartRepository
}
