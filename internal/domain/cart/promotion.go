package cart

import "math"

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

func (p *Promotion) CalculatePrice(price float64, qty int64) float64 {
	switch p.PromotionType {
	case PercentageDiscount:
		discountedPrice := price * (100 - float64(p.Discount)) / 100.0
		return discountedPrice * float64(qty)
	case Buy1Get1Free:
		paidQty := math.Ceil(float64(qty) / 2) // 3 / 2 = 1.5 := 2
		return price * paidQty
	}
	return 0
}
