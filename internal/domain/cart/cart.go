package cart

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type (
	CartItem struct {
		Product  Product
		Quantity int64
	}

	Cart struct {
		Items                  map[string]*CartItem
		Promotion              map[string]*Promotion
		TotalDiscountPromotion *Promotion
	}
)

func NewCart() *Cart {
	return &Cart{
		Items:     make(map[string]*CartItem),
		Promotion: make(map[string]*Promotion),
	}
}

func (c *Cart) AddProduct(product Product, quantity int64) error {
	if err := ValidateProduct(product); err != nil {
		return fmt.Errorf("invalid product: %w", err)
	}

	if err := ValidateQuantity(quantity); err != nil {
		return fmt.Errorf("invalid quantity: %w", err)
	}

	if item, ok := c.Items[product.ID]; ok {
		item.Quantity += quantity
		return nil
	}

	c.Items[product.ID] = &CartItem{Product: product, Quantity: quantity}
	return nil
}

func (c *Cart) AddPromotion(promotion Promotion) {
	if promotion.PromotionType == TotalDiscount {
		c.TotalDiscountPromotion = &promotion
		return
	}
	if _, ok := c.Promotion[promotion.ProductID]; ok {
		return
	}

	c.Promotion[promotion.ProductID] = &promotion
}

func (c *Cart) CalculateTotal() decimal.Decimal {
	total := decimal.Zero
	for _, item := range c.Items {
		// Apply product discount first
		discountedPrice := item.Product.GetDiscountedPrice()
		qty := item.Quantity
		promo, hasPromo := c.Promotion[item.Product.ID]

		if !hasPromo {
			qtyDecimal := decimal.NewFromInt(qty)
			total = total.Add(discountedPrice.Mul(qtyDecimal))
			continue
		}

		// Apply promotions to the discounted price
		total = total.Add(promo.CalculatePrice(discountedPrice, qty))
	}

	if c.TotalDiscountPromotion != nil {
		discount := decimal.NewFromInt(c.TotalDiscountPromotion.Discount)
		hundred := decimal.NewFromInt(100)
		total = total.Mul(hundred.Sub(discount)).Div(hundred)
	}

	return total
}

func DisplayPrice(price decimal.Decimal) string {
	// Convert decimal price to string with 2 decimal places
	return price.StringFixed(2)
}

// ValidateProduct validates product data
func ValidateProduct(product Product) error {
	if strings.TrimSpace(product.ID) == "" {
		return errors.New("product ID cannot be empty")
	}

	if product.Price.LessThan(decimal.Zero) {
		return errors.New("product price cannot be negative")
	}

	if !product.ValidateDiscount() {
		return errors.New("product discount must be between 0 and 100")
	}

	return nil
}

// ValidateQuantity validates quantity value
func ValidateQuantity(quantity int64) error {
	if quantity <= 0 {
		return errors.New("quantity must be a positive integer")
	}

	return nil
}
