package main

import (
	"fmt"
)

// 1. Service Definition
// แทนที่จะ Hardcode logic ไว้ใน Method, เราใช้ Field ที่เป็น Function แทน
type DiscountService struct {
	// calculateDiscount เป็น function field ที่สามารถเปลี่ยนไส้ในได้
	calculateDiscount func(amount float64) float64
}

// 2. Factory / Constructor
// รับ config เข้ามาเพื่อตัดสินใจว่าจะใช้ไส้ในแบบไหน
func NewDiscountService(promoMode string) *DiscountService {
	var strategy func(float64) float64

	// Logic การเลือก Strategy (Switch logic) อยู่ที่ตอนสร้าง Service
	switch promoMode {
	case "DOUBLE_DAY":
		strategy = superSaleDiscountFunc()
	case "VIP":
		strategy = vipDiscountFunc()
	default:
		strategy = standardDiscountFunc()
	}

	return &DiscountService{
		calculateDiscount: strategy,
	}
}

// 3. Public Method
// เรียกใช้ function ที่เก็บไว้ใน field
func (s *DiscountService) GetFinalPrice(price float64) float64 {
	discount := s.calculateDiscount(price)
	return price - discount
}

// ---------------------------------------------------------
// Implementation Strategies (Hidden/Private Functions)
// ---------------------------------------------------------

func standardDiscountFunc() func(float64) float64 {
	return func(amount float64) float64 {
		fmt.Println("🤖 Applying Standard Discount (5%)")
		return amount * 0.05
	}
}

func superSaleDiscountFunc() func(float64) float64 {
	return func(amount float64) float64 {
		fmt.Println("🔥 Applying 11.11 Super Sale Discount (50%)")
		return amount * 0.50
	}
}

func vipDiscountFunc() func(float64) float64 {
	return func(amount float64) float64 {
		fmt.Println("💎 Applying VIP Flat Discount (-100 THB)")
		if amount > 100 {
			return 100.0
		}
		return amount
	}
}

// ---------------------------------------------------------
// Demo Usage
// ---------------------------------------------------------

func main() {
	price := 1000.0

	fmt.Println("--- Scenario 1: Normal Day ---")
	// จำลองว่าอ่าน Config มาแล้วได้ค่า empty
	svcNormal := NewDiscountService("")
	fmt.Printf("Final Price: %.2f\n\n", svcNormal.GetFinalPrice(price))

	fmt.Println("--- Scenario 2: 11.11 Campaign ---")
	// จำลองว่าอ่าน Config มาได้ "DOUBLE_DAY"
	svcSale := NewDiscountService("DOUBLE_DAY")
	fmt.Printf("Final Price: %.2f\n\n", svcSale.GetFinalPrice(price))

	fmt.Println("--- Scenario 3: VIP Customer ---")
	svcVIP := NewDiscountService("VIP")
	fmt.Printf("Final Price: %.2f\n\n", svcVIP.GetFinalPrice(price))
}
