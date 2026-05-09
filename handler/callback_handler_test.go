package handler

import (
	"testing"

	"github.com/Rhymond/go-money"
)

func TestDecimalise(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		currency money.Currency
		want     int64
		wantErr  bool
	}{
		{"SGD 5.50", 5.50, *money.GetCurrency("SGD"), 550, false},
		{"SGD 100", 100, *money.GetCurrency("SGD"), 10000, false},
		{"SGD 0.01", 0.01, *money.GetCurrency("SGD"), 1, false},
		{"SGD 0", 0, *money.GetCurrency("SGD"), 0, false},
		{"JPY 100", 100, *money.GetCurrency("JPY"), 100, false},
		{"SGD -5.50", -5.50, *money.GetCurrency("SGD"), -550, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decimalise(tt.value, tt.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("decimalise(%v, %v) error = %v, wantErr %v", tt.value, tt.currency, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("decimalise(%v, %v) = %d, want %d", tt.value, tt.currency, got, tt.want)
			}
		})
	}
}

func TestRemoveNonNumeric(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"digits only", "12345", "12345"},
		{"with dots", "12.34", "1234"},
		{"with text", "abc123def", "123"},
		{"negative at start", "-100", "-100"},
		{"negative mid-string", "abc-100", "100"},
		{"empty", "", ""},
		{"all non-numeric", "abc", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeNonNumeric(tt.input)
			if got != tt.want {
				t.Errorf("removeNonNumeric(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
