package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/pkittipat/try-cart/internal/domain/cart"
	"github.com/pkittipat/try-cart/internal/domain/repository"
)

var (
	ErrCartNotFound    = errors.New("cart not found")
	ErrCartExists      = errors.New("cart already exists")
	ErrInvalidCartID   = errors.New("invalid cart ID")
	ErrInvalidUserID   = errors.New("invalid user ID")
)

type CartData struct {
	ID        string
	UserID    string
	Cart      *cart.Cart
	CreatedAt time.Time
	UpdatedAt time.Time
}

type cartRepository struct {
	mu    sync.RWMutex
	carts map[string]*CartData
	userCarts map[string]string // userID -> cartID mapping
}

func NewCartRepository() repository.Cart {
	return &cartRepository{
		carts:     make(map[string]*CartData),
		userCarts: make(map[string]string),
	}
}

func (r *cartRepository) Create(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", ErrInvalidUserID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user already has a cart
	if existingCartID, exists := r.userCarts[userID]; exists {
		return existingCartID, ErrCartExists
	}

	// Generate cart ID (simple implementation)
	cartID := fmt.Sprintf("cart_%s_%d", userID, time.Now().UnixNano())
	
	now := time.Now()
	cartData := &CartData{
		ID:        cartID,
		UserID:    userID,
		Cart:      cart.NewCart(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.carts[cartID] = cartData
	r.userCarts[userID] = cartID

	return cartID, nil
}

func (r *cartRepository) GetByID(ctx context.Context, cartID string) (*cart.Cart, error) {
	if cartID == "" {
		return nil, ErrInvalidCartID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	cartData, exists := r.carts[cartID]
	if !exists {
		return nil, ErrCartNotFound
	}

	return cartData.Cart, nil
}

func (r *cartRepository) GetByUserID(ctx context.Context, userID string) (*cart.Cart, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	cartID, exists := r.userCarts[userID]
	if !exists {
		return nil, ErrCartNotFound
	}

	cartData, exists := r.carts[cartID]
	if !exists {
		return nil, ErrCartNotFound
	}

	return cartData.Cart, nil
}

func (r *cartRepository) Update(ctx context.Context, cartID string, updatedCart *cart.Cart) error {
	if cartID == "" {
		return ErrInvalidCartID
	}
	if updatedCart == nil {
		return errors.New("cart cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	cartData, exists := r.carts[cartID]
	if !exists {
		return ErrCartNotFound
	}

	cartData.Cart = updatedCart
	cartData.UpdatedAt = time.Now()

	return nil
}

func (r *cartRepository) Delete(ctx context.Context, cartID string) error {
	if cartID == "" {
		return ErrInvalidCartID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	cartData, exists := r.carts[cartID]
	if !exists {
		return ErrCartNotFound
	}

	delete(r.carts, cartID)
	delete(r.userCarts, cartData.UserID)

	return nil
}

func (r *cartRepository) Exists(ctx context.Context, cartID string) (bool, error) {
	if cartID == "" {
		return false, ErrInvalidCartID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.carts[cartID]
	return exists, nil
}
