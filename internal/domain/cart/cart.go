package cart

import (
	"fmt"
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

func (c *Cart) AddProduct(product Product, quantity int64) {
	if item, ok := c.Items[product.ID]; ok {
		item.Quantity += quantity
		return
	}

	c.Items[product.ID] = &CartItem{Product: product, Quantity: quantity}
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
