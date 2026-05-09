package util

import (
	"testing"

	"github.com/Rhymond/go-money"
)

func TestGetFloatFormatter(t *testing.T) {
	tests := []struct {
		name    string
		curr    money.Currency
		want    string
	}{
		{"SGD (2 decimal places)", *money.GetCurrency("SGD"), "%.2f"},
		{"JPY (0 decimal places)", *money.GetCurrency("JPY"), "%f"},
		{"Fraction 3", money.Currency{Fraction: 3}, "%.3f"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFloatFormatter(tt.curr)
			if got != tt.want {
				t.Errorf("GetFloatFormatter() = %q, want %q", got, tt.want)
			}
		})
	}
}
