package handler

import (
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
)

type CartHandler struct {
	cartUseCase services.CartUseCase
}
