package cart

type (
	Product struct {
		ID          string
		Description string
		Price       float64 // as decimal price
		Discount    int64   // percentage discount (0-100)
	}
)

// GetDiscountedPrice calculates the final price after applying the product discount
func (p Product) GetDiscountedPrice() float64 {
	if p.Discount <= 0 || p.Discount > 100 {
		return p.Price
	}
	return p.Price * (100 - float64(p.Discount)) / 100.0
}

// ValidateDiscount checks if the discount percentage is valid
func (p Product) ValidateDiscount() bool {
	return p.Discount >= 0 && p.Discount <= 100
}
