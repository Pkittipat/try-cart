package cart

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCart_AddProduct(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		quantity int64
		setup    func(*Cart)
		want     map[string]*CartItem
	}{
		{
			name: "add new product",
			product: Product{
				ID:    "1",
				Price: decimal.NewFromFloat(10.00),
			},
			quantity: 2,
			setup:    func(c *Cart) {},
			want: map[string]*CartItem{
				"1": {
					Product: Product{
						ID:    "1",
						Price: decimal.NewFromFloat(10.00),
					},
					Quantity: 2,
				},
			},
		},
		{
			name: "add existing product",
			product: Product{
				ID:    "1",
				Price: decimal.NewFromFloat(10.00),
			},
			quantity: 3,
			setup: func(c *Cart) {
				err := c.AddProduct(Product{
					ID:    "1",
					Price: decimal.NewFromFloat(10.00),
				}, 2)
				assert.NoError(t, err)
			},
			want: map[string]*CartItem{
				"1": {
					Product: Product{
						ID:    "1",
						Price: decimal.NewFromFloat(10.00),
					},
					Quantity: 5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := NewCart()
			tt.setup(cart)
			err := cart.AddProduct(tt.product, tt.quantity)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, cart.Items)
		})
	}
}

func TestCart_AddPromotion(t *testing.T) {
	tests := []struct {
		name      string
		promotion Promotion
		setup     func(*Cart)
		want      map[string]*Promotion
	}{
		{
			name: "add new promotion",
			promotion: Promotion{
				ProductID:     "1",
				PromotionType: PercentageDiscount,
				Discount:      10,
			},
			setup: func(c *Cart) {},
			want: map[string]*Promotion{
				"1": {
					ProductID:     "1",
					PromotionType: PercentageDiscount,
					Discount:      10,
				},
			},
		},
		{
			name: "add existing promotion",
			promotion: Promotion{
				ProductID:     "1",
				PromotionType: PercentageDiscount,
				Discount:      20,
			},
			setup: func(c *Cart) {
				c.AddPromotion(Promotion{
					ProductID:     "1",
					PromotionType: PercentageDiscount,
					Discount:      10,
				})
			},
			want: map[string]*Promotion{
				"1": {
					ProductID:     "1",
					PromotionType: PercentageDiscount,
					Discount:      10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := NewCart()
			tt.setup(cart)
			cart.AddPromotion(tt.promotion)
			assert.Equal(t, tt.want, cart.Promotion)
		})
	}
}

func TestCart_CalculateTotal(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*Cart)
		want  decimal.Decimal
	}{
		{
			name:  "no items",
			setup: func(c *Cart) {},
			want:  decimal.NewFromFloat(0.0),
		},
		{
			name: "items without promotion",
			setup: func(c *Cart) {
				err := c.AddProduct(Product{
					ID:    "1",
					Price: decimal.NewFromFloat(10.00),
				}, 2)
				assert.NoError(t, err)
				err = c.AddProduct(Product{
					ID:    "2",
					Price: decimal.NewFromFloat(20.00),
				}, 1)
				assert.NoError(t, err)
			},
			want: decimal.NewFromFloat(40.00), // (10.00 * 2) + (20.00 * 1)
		},
		{
			name: "items with percentage discount",
			setup: func(c *Cart) {
				err := c.AddProduct(Product{
					ID:    "1",
					Price: decimal.NewFromFloat(10.00),
				}, 2)
				assert.NoError(t, err)
				c.AddPromotion(Promotion{
					ProductID:     "1",
					PromotionType: PercentageDiscount,
					Discount:      10,
				})
			},
			want: decimal.NewFromFloat(18.00), // (10.00 * 0.9 * 2)
		},
		{
			name: "items with buy 1 get 1 free",
			setup: func(c *Cart) {
				err := c.AddProduct(Product{
					ID:    "1",
					Price: decimal.NewFromFloat(10.00),
				}, 3)
				assert.NoError(t, err)
				c.AddPromotion(Promotion{
					ProductID:     "1",
					PromotionType: Buy1Get1Free,
				})
			},
			want: decimal.NewFromFloat(20.00), // (10.00 * 2) - pay for 2 items, get 1 free
		},
		{
			name: "items with total discount",
			setup: func(c *Cart) {
				err := c.AddProduct(Product{
					ID:    "1",
					Price: decimal.NewFromFloat(10.00),
				}, 2)
				assert.NoError(t, err)
				err = c.AddProduct(Product{
					ID:    "2",
					Price: decimal.NewFromFloat(20.00),
				}, 1)
				assert.NoError(t, err)
				c.AddPromotion(Promotion{
					PromotionType: TotalDiscount,
					Discount:      10,
				})
			},
			want: decimal.NewFromFloat(36.00), // (40.00 * 0.9)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := NewCart()
			tt.setup(cart)
			got := cart.CalculateTotal()
			assert.True(t, tt.want.Equal(got), "Expected %s, got %s", tt.want.String(), got.String())
		})
	}
}

func TestDisplayPrice(t *testing.T) {
	tests := []struct {
		name  string
		price decimal.Decimal
		want  string
	}{
		{
			name:  "positive price",
			price: decimal.NewFromFloat(123.45),
			want:  "123.45",
		},
		{
			name:  "zero price",
			price: decimal.NewFromFloat(0.00),
			want:  "0.00",
		},
		{
			name:  "price with two decimal places",
			price: decimal.NewFromFloat(100.00),
			want:  "100.00",
		},
		{
			name:  "price with one decimal place",
			price: decimal.NewFromFloat(123.40),
			want:  "123.40",
		},
		{
			name:  "precise decimal from string",
			price: decimal.RequireFromString("9999999999999999.99"),
			want:  "9999999999999999.99", // Precise decimal maintains accuracy
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DisplayPrice(tt.price)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProduct_GetDiscountedPrice(t *testing.T) {
	tests := []struct {
		name    string
		product Product
		want    decimal.Decimal
	}{
		{
			name: "no discount",
			product: Product{
				ID:       "1",
				Price:    decimal.NewFromFloat(100.00),
				Discount: 0,
			},
			want: decimal.NewFromFloat(100.00),
		},
		{
			name: "10% discount",
			product: Product{
				ID:       "1",
				Price:    decimal.NewFromFloat(100.00),
				Discount: 10,
			},
			want: decimal.NewFromFloat(90.00),
		},
		{
			name: "50% discount",
			product: Product{
				ID:       "1",
				Price:    decimal.NewFromFloat(100.00),
				Discount: 50,
			},
			want: decimal.NewFromFloat(50.00),
		},
		{
			name: "100% discount",
			product: Product{
				ID:       "1",
				Price:    decimal.NewFromFloat(100.00),
				Discount: 100,
			},
			want: decimal.NewFromFloat(0.00),
		},
		{
			name: "invalid discount over 100%",
			product: Product{
				ID:       "1",
				Price:    decimal.NewFromFloat(100.00),
				Discount: 150,
			},
			want: decimal.NewFromFloat(100.00),
		},
		{
			name: "negative discount",
			product: Product{
				ID:       "1",
				Price:    decimal.NewFromFloat(100.00),
				Discount: -10,
			},
			want: decimal.NewFromFloat(100.00),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.product.GetDiscountedPrice()
			assert.True(t, tt.want.Equal(got), "Expected %s, got %s", tt.want.String(), got.String())
		})
	}
}

func TestProduct_ValidateDiscount(t *testing.T) {
	tests := []struct {
		name     string
		discount int64
		want     bool
	}{
		{
			name:     "valid discount 0%",
			discount: 0,
			want:     true,
		},
		{
			name:     "valid discount 50%",
			discount: 50,
			want:     true,
		},
		{
			name:     "valid discount 100%",
			discount: 100,
			want:     true,
		},
		{
			name:     "invalid discount over 100%",
			discount: 101,
			want:     false,
		},
		{
			name:     "invalid negative discount",
			discount: -1,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := Product{
				ID:       "1",
				Price:    decimal.NewFromFloat(100.00),
				Discount: tt.discount,
			}
			got := product.ValidateDiscount()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCart_CalculateTotal_WithProductDiscounts(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*Cart)
		want  decimal.Decimal
	}{
		{
			name: "single product with discount",
			setup: func(cart *Cart) {
				product := Product{
					ID:       "A",
					Price:    decimal.NewFromFloat(100.00), // 100.00
					Discount: 10,     // 10% discount
				}
				err := cart.AddProduct(product, 1)
				assert.NoError(t, err)
			},
			want: decimal.NewFromFloat(90.00), // 90.00
		},
		{
			name: "multiple products with different discounts",
			setup: func(cart *Cart) {
				productA := Product{
					ID:       "A",
					Price:    decimal.NewFromFloat(100.00), // 100.00
					Discount: 10,     // 10% discount
				}
				productB := Product{
					ID:       "B",
					Price:    decimal.NewFromFloat(50.00), // 50.00
					Discount: 20,    // 20% discount
				}
				err := cart.AddProduct(productA, 1) // 90.00
				assert.NoError(t, err)
				err = cart.AddProduct(productB, 2) // 40.00 * 2 = 80.00
				assert.NoError(t, err)
			},
			want: decimal.NewFromFloat(170.00), // 90.00 + 80.00 = 170.00
		},
		{
			name: "product discount with promotion",
			setup: func(cart *Cart) {
				product := Product{
					ID:       "A",
					Price:    decimal.NewFromFloat(100.00), // 100.00
					Discount: 10,     // 10% discount -> 90.00
				}
				err := cart.AddProduct(product, 2)
				assert.NoError(t, err)
				// Apply 18% promotion discount on already discounted price
				cart.AddPromotion(Promotion{
					ProductID:     "A",
					PromotionType: PercentageDiscount,
					Discount:      18,
				})
			},
			want: decimal.NewFromFloat(147.60), // (90.00 * 2) * 0.82 = 147.60
		},
		{
			name: "product discount with buy 1 get 1 free",
			setup: func(cart *Cart) {
				product := Product{
					ID:       "A",
					Price:    decimal.NewFromFloat(100.00), // 100.00
					Discount: 20,     // 20% discount -> 80.00
				}
				err := cart.AddProduct(product, 3)
				assert.NoError(t, err)
				cart.AddPromotion(Promotion{
					ProductID:     "A",
					PromotionType: Buy1Get1Free,
				})
			},
			want: decimal.NewFromFloat(160.00), // ceil(3/2) * 80.00 = 2 * 80.00 = 160.00
		},
		{
			name: "product discount with total discount",
			setup: func(cart *Cart) {
				product := Product{
					ID:       "A",
					Price:    decimal.NewFromFloat(100.00), // 100.00
					Discount: 10,     // 10% discount -> 90.00
				}
				err := cart.AddProduct(product, 2)
				assert.NoError(t, err)
				cart.AddPromotion(Promotion{
					PromotionType: TotalDiscount,
					Discount:      15, // 15% total discount
				})
			},
			want: decimal.NewFromFloat(153.00), // (90.00 * 2) * 0.85 = 153.00
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := NewCart()
			tt.setup(cart)
			got := cart.CalculateTotal()
			assert.True(t, tt.want.Equal(got), "Expected %s, got %s", tt.want.String(), got.String())
		})
	}
}

func TestCart_AddProduct_Validation(t *testing.T) {
	tests := []struct {
		name        string
		product     Product
		quantity    int64
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid product and quantity",
			product: Product{
				ID:    "1",
				Price: decimal.NewFromFloat(10.00),
			},
			quantity:    2,
			expectError: false,
		},
		{
			name: "empty product ID",
			product: Product{
				ID:    "",
				Price: decimal.NewFromFloat(10.00),
			},
			quantity:    2,
			expectError: true,
			errorMsg:    "invalid product: product ID cannot be empty",
		},
		{
			name: "whitespace product ID",
			product: Product{
				ID:    "   ",
				Price: decimal.NewFromFloat(10.00),
			},
			quantity:    2,
			expectError: true,
			errorMsg:    "invalid product: product ID cannot be empty",
		},
		{
			name: "negative product price",
			product: Product{
				ID:    "1",
				Price: decimal.NewFromFloat(-10.00),
			},
			quantity:    2,
			expectError: true,
			errorMsg:    "invalid product: product price cannot be negative",
		},
		{
			name: "invalid product discount over 100",
			product: Product{
				ID:       "1",
				Price:    decimal.NewFromFloat(10.00),
				Discount: 150,
			},
			quantity:    2,
			expectError: true,
			errorMsg:    "invalid product: product discount must be between 0 and 100",
		},
		{
			name: "invalid product discount negative",
			product: Product{
				ID:       "1",
				Price:    decimal.NewFromFloat(10.00),
				Discount: -10,
			},
			quantity:    2,
			expectError: true,
			errorMsg:    "invalid product: product discount must be between 0 and 100",
		},
		{
			name: "zero quantity",
			product: Product{
				ID:    "1",
				Price: decimal.NewFromFloat(10.00),
			},
			quantity:    0,
			expectError: true,
			errorMsg:    "invalid quantity: quantity must be a positive integer",
		},
		{
			name: "negative quantity",
			product: Product{
				ID:    "1",
				Price: decimal.NewFromFloat(10.00),
			},
			quantity:    -5,
			expectError: true,
			errorMsg:    "invalid quantity: quantity must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := NewCart()
			err := cart.AddProduct(tt.product, tt.quantity)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateProduct(t *testing.T) {
	tests := []struct {
		name        string
		product     Product
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid product",
			product: Product{
				ID:       "valid-id",
				Price:    decimal.NewFromFloat(100.00),
				Discount: 10,
			},
			expectError: false,
		},
		{
			name: "empty product ID",
			product: Product{
				ID:    "",
				Price: decimal.NewFromFloat(100.00),
			},
			expectError: true,
			errorMsg:    "product ID cannot be empty",
		},
		{
			name: "whitespace product ID",
			product: Product{
				ID:    "   ",
				Price: decimal.NewFromFloat(100.00),
			},
			expectError: true,
			errorMsg:    "product ID cannot be empty",
		},
		{
			name: "negative price",
			product: Product{
				ID:    "valid-id",
				Price: decimal.NewFromFloat(-50.00),
			},
			expectError: true,
			errorMsg:    "product price cannot be negative",
		},
		{
			name: "invalid discount over 100",
			product: Product{
				ID:       "valid-id",
				Price:    decimal.NewFromFloat(100.00),
				Discount: 150,
			},
			expectError: true,
			errorMsg:    "product discount must be between 0 and 100",
		},
		{
			name: "invalid negative discount",
			product: Product{
				ID:       "valid-id",
				Price:    decimal.NewFromFloat(100.00),
				Discount: -10,
			},
			expectError: true,
			errorMsg:    "product discount must be between 0 and 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProduct(tt.product)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateQuantity(t *testing.T) {
	tests := []struct {
		name        string
		quantity    int64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid positive quantity",
			quantity:    5,
			expectError: false,
		},
		{
			name:        "valid quantity 1",
			quantity:    1,
			expectError: false,
		},
		{
			name:        "zero quantity",
			quantity:    0,
			expectError: true,
			errorMsg:    "quantity must be a positive integer",
		},
		{
			name:        "negative quantity",
			quantity:    -3,
			expectError: true,
			errorMsg:    "quantity must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuantity(tt.quantity)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
