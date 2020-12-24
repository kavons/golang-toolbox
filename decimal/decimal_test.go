package decimal

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func TestDecimalPackage(t *testing.T) {
	price, err := decimal.NewFromString("136.02")
	if err != nil {
		panic(err)
	}

	quantity := decimal.NewFromFloat(3)

	fee, _ := decimal.NewFromString(".035")
	taxRate, _ := decimal.NewFromString(".08875")

	subtotal := price.Mul(quantity)

	preTax := subtotal.Mul(fee.Add(decimal.NewFromFloat(1)))

	total := preTax.Mul(taxRate.Add(decimal.NewFromFloat(1)))

	fmt.Println("Subtotal:", subtotal)                      // Subtotal: 408.06
	fmt.Println("Pre-tax:", preTax)                         // Pre-tax: 422.3421
	fmt.Println("Taxes:", total.Sub(preTax))                // Taxes: 37.482861375
	fmt.Println("Total:", total)                            // Total: 459.824961375
	fmt.Println("Tax rate:", total.Sub(preTax).Div(preTax)) // Tax rate: 0.08875

	a := 8.842700
	b := 8.8427001
	fmt.Println(decimal.NewFromFloat(a).Cmp(decimal.NewFromFloat(b)))
	fmt.Println(a == b)
}

func TestDecimalCeil(t *testing.T) {
	volume := 9.000000
	fee := 0.0005

	// method one
	r := decimal.NewFromFloat(volume).Mul(decimal.NewFromFloat(fee))
	s := r.String()
	f, _ := r.Float64()
	fmt.Println(s, f)

	// method two
	r = decimal.NewFromFloat(volume).Mul(decimal.NewFromFloat(fee)).Mul(decimal.NewFromFloat(10000)).Ceil().Div(decimal.NewFromFloat(10000))
	s = r.String()
	f, _ = r.Float64()
	fmt.Println(s, f)
}
