package usecase

import (
	"context"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
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

func (c *cartUseCase) RemoveFromCart(ctx context.Context, productDetailsID int, userId int) error {
	err := c.cartRepo.RemoveFromCart(ctx, productDetailsID, userId)
	return err

}

func (c *cartUseCase) ViewCart(ctx context.Context, userId int) (response.ViewCart, error) {
	viewCart, err := c.cartRepo.ViewCart(ctx, userId)
	return viewCart, err

}

func (c *cartUseCase) FindCartItemsByUserID(ctx context.Context, user_id int) ([]domain.CartItems, error) {
	cartItems, err := c.cartRepo.FindCartItemsByUserID(ctx, user_id)
	if err != nil {
		return []domain.CartItems{}, err
	}
	return cartItems, nil
}
