package repository

import (
	"context"

	"github.com/pkittipat/try-cart/internal/domain/cart"
)

type Cart interface {
	// Create creates a new cart and returns the cart ID
	Create(ctx context.Context, userID string) (string, error)
	
	// GetByID retrieves a cart by its ID
	GetByID(ctx context.Context, cartID string) (*cart.Cart, error)
	
	// GetByUserID retrieves a cart by user ID
	GetByUserID(ctx context.Context, userID string) (*cart.Cart, error)
	
	// Update updates an existing cart
	Update(ctx context.Context, cartID string, cart *cart.Cart) error
	
	// Delete removes a cart by its ID
	Delete(ctx context.Context, cartID string) error
	
	// Exists checks if a cart exists by ID
	Exists(ctx context.Context, cartID string) (bool, error)
}
