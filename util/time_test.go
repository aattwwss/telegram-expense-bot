package util

import (
	"testing"
	"time"
)

func TestParseMonthFromString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Month
	}{
		{"numeric jan", "1", time.January},
		{"numeric dec", "12", time.December},
		{"short jan", "jan", time.January},
		{"short feb", "feb", time.February},
		{"short mar", "mar", time.March},
		{"short apr", "apr", time.April},
		{"short may", "may", time.May},
		{"short jun", "jun", time.June},
		{"short jul", "jul", time.July},
		{"short aug", "aug", time.August},
		{"short sep", "sep", time.September},
		{"short oct", "oct", time.October},
		{"short nov", "nov", time.November},
		{"short dec", "dec", time.December},
		{"full january", "january", time.January},
		{"full February", "February", time.February},
		{"case insensitive", "JANUARY", time.January},
		{"unknown returns current", "xyz", time.Now().Month()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseMonthFromString(tt.input)
			if got != tt.want {
				t.Errorf("ParseMonthFromString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseMonthYearFromMessage(t *testing.T) {
	now := time.Now()
	currentMonth := now.Month()
	currentYear := now.Year()

	tests := []struct {
		name      string
		input     string
		wantMonth time.Month
		wantYear  int
	}{
		{"no args returns current", "/list", currentMonth, currentYear},
		{"month only", "/list 3", time.March, currentYear},
		{"short month name", "/list Jan", time.January, currentYear},
		{"month and year", "/list 2 2022", time.February, 2022},
		{"month name and year", "/list February 2022", time.February, 2022},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMonth, gotYear := ParseMonthYearFromMessage(tt.input)
			if gotMonth != tt.wantMonth || gotYear != tt.wantYear {
				t.Errorf("ParseMonthYearFromMessage(%q) = (%v, %d), want (%v, %d)",
					tt.input, gotMonth, gotYear, tt.wantMonth, tt.wantYear)
			}
		})
	}
}

func TestYearMonthString(t *testing.T) {
	tests := []struct {
		name    string
		ym      YearMonth
		layout  string
		want    string
		wantErr bool
	}{
		{"standard layout", YearMonth{time.January, 2023}, "2006-01", "2023-01", false},
		{"december", YearMonth{time.December, 2023}, "2006-01", "2023-12", false},
		{"invalid layout", YearMonth{time.January, 2023}, "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ym.String(tt.layout)
			if (err != nil) != tt.wantErr {
				t.Errorf("YearMonth.String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("YearMonth.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
