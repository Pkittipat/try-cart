package main

import (
	"fmt"

	"github.com/pkittipat/try-cart/internal/domain/cart"
)

func main() {
	fmt.Println("=== Shopping Cart with Product Discounts Demo ===")

	// Create products with different discount levels
	productA := cart.Product{ID: "A", Price: 100.00, Discount: 10} // 100.00 with 10% discount	
	productB := cart.Product{ID: "B", Price: 200.00, Discount: 0}  // 200.00 no discount  	
	productC := cart.Product{ID: "C", Price: 50.53, Discount: 15}  // 50.53 with 15% discount

	fmt.Printf("Products:\n")
	fmt.Printf("Product A: Original %s, Discount %d%%, Final %s\n", 
		cart.DisplayPrice(productA.Price), productA.Discount, cart.DisplayPrice(productA.GetDiscountedPrice()))
	fmt.Printf("Product B: Original %s, Discount %d%%, Final %s\n", 
		cart.DisplayPrice(productB.Price), productB.Discount, cart.DisplayPrice(productB.GetDiscountedPrice()))
	fmt.Printf("Product C: Original %s, Discount %d%%, Final %s\n", 
		cart.DisplayPrice(productC.Price), productC.Discount, cart.DisplayPrice(productC.GetDiscountedPrice()))
	fmt.Println()

	// Create shopping cart and add products
	shoppingCart := cart.NewCart()
	shoppingCart.AddProduct(productA, 3)
	shoppingCart.AddProduct(productB, 1) 
	shoppingCart.AddProduct(productC, 2)

	fmt.Printf("Cart contents:\n")
	fmt.Printf("- Product A x3: %s\n", cart.DisplayPrice(productA.GetDiscountedPrice()*3))
	fmt.Printf("- Product B x1: %s\n", cart.DisplayPrice(productB.GetDiscountedPrice()*1))
	fmt.Printf("- Product C x2: %s\n", cart.DisplayPrice(productC.GetDiscountedPrice()*2))
	fmt.Printf("Subtotal (with product discounts): %s\n", cart.DisplayPrice(shoppingCart.CalculateTotal()))
	fmt.Println()

	// Apply additional promotions on top of product discounts
	fmt.Println("Applying additional promotions:")
	shoppingCart.AddPromotion(cart.Promotion{ProductID: "A", PromotionType: cart.PercentageDiscount, Discount: 18})
	fmt.Printf("- Product A: Additional 18%% promotion discount\n")
	
	shoppingCart.AddPromotion(cart.Promotion{ProductID: "C", PromotionType: cart.Buy1Get1Free})
	fmt.Printf("- Product C: Buy 1 Get 1 Free promotion\n")
	
	fmt.Printf("Final Total: %s\n", cart.DisplayPrice(shoppingCart.CalculateTotal()))
	fmt.Println()

	// Demonstrate validation
	fmt.Println("=== Product Discount Validation ===")
	invalidProduct := cart.Product{ID: "INVALID", Price: 100.00, Discount: 150}
	fmt.Printf("Product with 150%% discount is valid: %t\n", invalidProduct.ValidateDiscount())
	fmt.Printf("Invalid product discounted price: %s (should be original price)\n", 
		cart.DisplayPrice(invalidProduct.GetDiscountedPrice()))
}
