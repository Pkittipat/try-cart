package cart

import "github.com/shopspring/decimal"

type (
	Product struct {
		ID          string
		Description string
		Price       decimal.Decimal // as decimal price
		Discount    int64           // percentage discount (0-100)
	}
)

// GetDiscountedPrice calculates the final price after applying the product discount
func (p Product) GetDiscountedPrice() decimal.Decimal {
	if p.Discount <= 0 || p.Discount > 100 {
		return p.Price
	}
	discount := decimal.NewFromInt(p.Discount)
	hundred := decimal.NewFromInt(100)
	return p.Price.Mul(hundred.Sub(discount)).Div(hundred)
}

// ValidateDiscount checks if the discount percentage is valid
func (p Product) ValidateDiscount() bool {
	return p.Discount >= 0 && p.Discount <= 100
}
