package service

import "github.com/pkittipat/try-cart/internal/domain/repository"

type CartService struct {
	cartRepo repository.Cart
}

func NewCartService(
	cartRepo repository.Cart,
) *CartService {
	return &CartService{
		cartRepo: cartRepo,
	}
}
