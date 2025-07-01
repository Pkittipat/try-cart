package cart

import (
	"errors"
	"fmt"
	"strings"
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

func (c *Cart) CalculateTotal() float64 {
	total := 0.0
	for _, item := range c.Items {
		// Apply product discount first
		discountedPrice := item.Product.GetDiscountedPrice()
		qty := item.Quantity
		promo, hasPromo := c.Promotion[item.Product.ID]

		if !hasPromo {
			total += (discountedPrice * float64(qty))
			continue
		}

		// Apply promotions to the discounted price
		total += promo.CalculatePrice(discountedPrice, qty)
	}

	if c.TotalDiscountPromotion != nil {
		discount := c.TotalDiscountPromotion.Discount
		total = total * (100 - float64(discount)) / 100.0
	}

	return total
}

func DisplayPrice(price float64) string {
	// Convert float64 price to string with 2 decimal places
	return fmt.Sprintf("%.2f", price)
}

// ValidateProduct validates product data
func ValidateProduct(product Product) error {
	if strings.TrimSpace(product.ID) == "" {
		return errors.New("product ID cannot be empty")
	}

	if product.Price < 0 {
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
