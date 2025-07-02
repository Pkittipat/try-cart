package cart

import (
	"math"

	"github.com/shopspring/decimal"
)

type PromotionType string

const (
	PercentageDiscount PromotionType = "percentageDiscount"
	Buy1Get1Free       PromotionType = "buy1Get1Free"
	TotalDiscount      PromotionType = "totalDiscount"
)

type (
	Promotion struct {
		ID            string
		Discount      int64
		ProductID     string
		PromotionType PromotionType
	}
)

func (p *Promotion) CalculatePrice(price decimal.Decimal, qty int64) decimal.Decimal {
	switch p.PromotionType {
	case PercentageDiscount:
		discount := decimal.NewFromInt(p.Discount)
		hundred := decimal.NewFromInt(100)
		quantity := decimal.NewFromInt(qty)
		discountedPrice := price.Mul(hundred.Sub(discount)).Div(hundred)
		return discountedPrice.Mul(quantity)
	case Buy1Get1Free:
		paidQty := math.Ceil(float64(qty) / 2) // 3 / 2 = 1.5 := 2
		paidQtyDecimal := decimal.NewFromFloat(paidQty)
		return price.Mul(paidQtyDecimal)
	}
	return decimal.Zero
}
