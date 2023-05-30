package usecase

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
)

type cartUseCase struct {
	cartRepo interfaces.CartRepository
}

func NewCartUseCase(cartRepo interfaces.CartRepository) services.CartUseCase {
	return &cartUseCase{
		cartRepo: cartRepo,
	}
}

func (c *cartUseCase) AddToCart(ctx context.Context, productDetailsID int, userID int) (domain.CartItems, error) {
	addedProduct, err := c.cartRepo.AddToCart(ctx, productDetailsID, userID)
	return addedProduct, err
}
