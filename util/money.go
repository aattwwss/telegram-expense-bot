package util

import (
	"fmt"
	"github.com/Rhymond/go-money"
)

func GetFloatFormatter(currency money.Currency) string {
	if currency.Fraction == 0 {
		return "%f"
	}
	return fmt.Sprintf("%%.%vf", currency.Fraction)
}
