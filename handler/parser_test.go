package handler

import (
	"testing"
)

func TestParseFloatStringFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"integer", "100 Groceries", "100", false},
		{"decimal", "5.50 Chicken Rice", "5.50", false},
		{"decimal with single digit cents", "5.5 Coffee", "5.5", false},
		{"negative", "-10 Refund", "-10", false},
		{"no cents", "100.", "100.", false},
		{"no float", "Chicken Rice", "", true},
		{"amount at end", "Chicken Rice 5.50", "", true},
		{"multiple dots", "5.5.0 test", "5.5", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFloatStringFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFloatStringFromString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseFloatStringFromString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
