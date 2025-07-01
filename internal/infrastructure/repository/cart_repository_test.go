package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkittipat/try-cart/internal/domain/cart"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartRepository_Create(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		setupFunc   func(*cartRepository)
		wantErr     error
		wantCartID  bool
	}{
		{
			name:       "successful creation",
			userID:     "user123",
			wantErr:    nil,
			wantCartID: true,
		},
		{
			name:       "empty userID",
			userID:     "",
			wantErr:    ErrInvalidUserID,
			wantCartID: false,
		},
		{
			name:   "user already has cart",
			userID: "existing_user",
			setupFunc: func(r *cartRepository) {
				r.userCarts["existing_user"] = "existing_cart_id"
			},
			wantErr:    ErrCartExists,
			wantCartID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewCartRepository().(*cartRepository)
			ctx := context.Background()

			if tt.setupFunc != nil {
				tt.setupFunc(repo)
			}

			cartID, err := repo.Create(ctx, tt.userID)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}

			if tt.wantCartID {
				assert.NotEmpty(t, cartID)
			} else {
				assert.Empty(t, cartID)
			}
		})
	}
}

func TestCartRepository_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		cartID    string
		setupFunc func(*cartRepository)
		wantErr   error
		wantCart  bool
	}{
		{
			name:   "successful retrieval",
			cartID: "test_cart_id",
			setupFunc: func(r *cartRepository) {
				r.carts["test_cart_id"] = &CartData{
					ID:     "test_cart_id",
					UserID: "user123",
					Cart:   cart.NewCart(),
				}
			},
			wantErr:  nil,
			wantCart: true,
		},
		{
			name:     "empty cartID",
			cartID:   "",
			wantErr:  ErrInvalidCartID,
			wantCart: false,
		},
		{
			name:     "non-existent cart",
			cartID:   "non_existent",
			wantErr:  ErrCartNotFound,
			wantCart: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewCartRepository().(*cartRepository)
			ctx := context.Background()

			if tt.setupFunc != nil {
				tt.setupFunc(repo)
			}

			result, err := repo.GetByID(ctx, tt.cartID)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}

			if tt.wantCart {
				assert.NotNil(t, result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestCartRepository_GetByUserID(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		setupFunc func(*cartRepository)
		wantErr   error
		wantCart  bool
	}{
		{
			name:   "successful retrieval",
			userID: "user123",
			setupFunc: func(r *cartRepository) {
				r.userCarts["user123"] = "cart_id_123"
				r.carts["cart_id_123"] = &CartData{
					ID:     "cart_id_123",
					UserID: "user123",
					Cart:   cart.NewCart(),
				}
			},
			wantErr:  nil,
			wantCart: true,
		},
		{
			name:     "empty userID",
			userID:   "",
			wantErr:  ErrInvalidUserID,
			wantCart: false,
		},
		{
			name:     "user not found",
			userID:   "non_existent_user",
			wantErr:  ErrCartNotFound,
			wantCart: false,
		},
		{
			name:   "user mapping exists but cart missing",
			userID: "orphaned_user",
			setupFunc: func(r *cartRepository) {
				r.userCarts["orphaned_user"] = "missing_cart_id"
			},
			wantErr:  ErrCartNotFound,
			wantCart: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewCartRepository().(*cartRepository)
			ctx := context.Background()

			if tt.setupFunc != nil {
				tt.setupFunc(repo)
			}

			result, err := repo.GetByUserID(ctx, tt.userID)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}

			if tt.wantCart {
				assert.NotNil(t, result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestCartRepository_Update(t *testing.T) {
	tests := []struct {
		name        string
		cartID      string
		updatedCart *cart.Cart
		setupFunc   func(*cartRepository)
		wantErr     string
	}{
		{
			name:   "successful update",
			cartID: "test_cart_id",
			updatedCart: func() *cart.Cart {
				c := cart.NewCart()
				c.AddProduct(cart.Product{ID: "A", Price: 10.00}, 2)
				return c
			}(),
			setupFunc: func(r *cartRepository) {
				r.carts["test_cart_id"] = &CartData{
					ID:     "test_cart_id",
					UserID: "user123",
					Cart:   cart.NewCart(),
				}
			},
			wantErr: "",
		},
		{
			name:        "empty cartID",
			cartID:      "",
			updatedCart: cart.NewCart(),
			wantErr:     ErrInvalidCartID.Error(),
		},
		{
			name:        "nil cart",
			cartID:      "test_cart_id",
			updatedCart: nil,
			wantErr:     "cart cannot be nil",
		},
		{
			name:        "non-existent cart",
			cartID:      "non_existent",
			updatedCart: cart.NewCart(),
			wantErr:     ErrCartNotFound.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewCartRepository().(*cartRepository)
			ctx := context.Background()

			if tt.setupFunc != nil {
				tt.setupFunc(repo)
			}

			err := repo.Update(ctx, tt.cartID, tt.updatedCart)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)

				// Verify the update took effect
				result, err := repo.GetByID(ctx, tt.cartID)
				require.NoError(t, err)
				assert.Equal(t, tt.updatedCart, result)
			}
		})
	}
}

func TestCartRepository_Delete(t *testing.T) {
	tests := []struct {
		name      string
		cartID    string
		setupFunc func(*cartRepository)
		wantErr   error
	}{
		{
			name:   "successful deletion",
			cartID: "test_cart_id",
			setupFunc: func(r *cartRepository) {
				r.carts["test_cart_id"] = &CartData{
					ID:     "test_cart_id",
					UserID: "user123",
					Cart:   cart.NewCart(),
				}
				r.userCarts["user123"] = "test_cart_id"
			},
			wantErr: nil,
		},
		{
			name:    "empty cartID",
			cartID:  "",
			wantErr: ErrInvalidCartID,
		},
		{
			name:    "non-existent cart",
			cartID:  "non_existent",
			wantErr: ErrCartNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewCartRepository().(*cartRepository)
			ctx := context.Background()

			if tt.setupFunc != nil {
				tt.setupFunc(repo)
			}

			err := repo.Delete(ctx, tt.cartID)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)

				// Verify deletion
				_, err = repo.GetByID(ctx, tt.cartID)
				assert.Equal(t, ErrCartNotFound, err)
			}
		})
	}
}

func TestCartRepository_Exists(t *testing.T) {
	tests := []struct {
		name       string
		cartID     string
		setupFunc  func(*cartRepository)
		wantExists bool
		wantErr    error
	}{
		{
			name:   "existing cart",
			cartID: "test_cart_id",
			setupFunc: func(r *cartRepository) {
				r.carts["test_cart_id"] = &CartData{
					ID:     "test_cart_id",
					UserID: "user123",
					Cart:   cart.NewCart(),
				}
			},
			wantExists: true,
			wantErr:    nil,
		},
		{
			name:       "non-existent cart",
			cartID:     "non_existent",
			wantExists: false,
			wantErr:    nil,
		},
		{
			name:       "empty cartID",
			cartID:     "",
			wantExists: false,
			wantErr:    ErrInvalidCartID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewCartRepository().(*cartRepository)
			ctx := context.Background()

			if tt.setupFunc != nil {
				tt.setupFunc(repo)
			}

			exists, err := repo.Exists(ctx, tt.cartID)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantExists, exists)
		})
	}
}

func TestCartRepository_ThreadSafety(t *testing.T) {
	repo := NewCartRepository()
	ctx := context.Background()

	const numGoroutines = 10
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			userID := fmt.Sprintf("user%d", id)
			cartID, err := repo.Create(ctx, userID)
			if err != nil {
				results <- err
				return
			}

			_, err = repo.GetByID(ctx, cartID)
			if err != nil {
				results <- err
				return
			}

			updatedCart := cart.NewCart()
			product := cart.Product{ID: fmt.Sprintf("product%d", id), Price: float64(id * 100)}
			updatedCart.AddProduct(product, int64(id))

			err = repo.Update(ctx, cartID, updatedCart)
			results <- err
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}